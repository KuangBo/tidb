// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package expression

import (
	"encoding/hex"
	"math/rand"
	"testing"

	. "github.com/pingcap/check"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/tidb/tablecodec"
	"github.com/pingcap/tidb/types"
)

type tidbKeyGener struct {
	inner *defaultGener
}

func (g *tidbKeyGener) gen() interface{} {
	tableID := g.inner.gen().(int64)
	var result []byte
	if rand.Intn(2) == 1 {
		// Generate a record key
		handle := g.inner.gen().(int64)
		result = tablecodec.EncodeRowKeyWithHandle(tableID, handle)
	} else {
		// Generate an index key
		idx := g.inner.gen().(int64)
		result = tablecodec.EncodeTableIndexPrefix(tableID, idx)
	}
	return hex.EncodeToString(result)
}

var vecBuiltinInfoCases = map[string][]vecExprBenchCase{
	ast.Version: {
		{retEvalType: types.ETString, childrenTypes: []types.EvalType{}},
	},
	ast.TiDBVersion: {
		{retEvalType: types.ETString, childrenTypes: []types.EvalType{}},
	},
	ast.CurrentUser: {},
	ast.FoundRows:   {},
	ast.Database: {
		{retEvalType: types.ETString, childrenTypes: []types.EvalType{}},
	},
	ast.User: {},
	ast.TiDBDecodeKey: {
		{
			retEvalType:   types.ETString,
			childrenTypes: []types.EvalType{types.ETString},
			geners: []dataGenerator{&tidbKeyGener{
				inner: &defaultGener{
					nullRation: 0,
					eType:      types.ETInt,
				},
			}},
		},
	},
	ast.RowCount: {
		{retEvalType: types.ETInt, childrenTypes: []types.EvalType{}},
	},
	ast.CurrentRole: {},
	ast.TiDBIsDDLOwner: {
		{retEvalType: types.ETInt, childrenTypes: []types.EvalType{}},
	},
	ast.ConnectionID: {},
	ast.LastInsertId: {
		{retEvalType: types.ETInt, childrenTypes: []types.EvalType{}},
		{retEvalType: types.ETInt, childrenTypes: []types.EvalType{types.ETInt}},
	},
}

func (s *testEvaluatorSuite) TestVectorizedBuiltinInfoFunc(c *C) {
	testVectorizedBuiltinFunc(c, vecBuiltinInfoCases)
}

func BenchmarkVectorizedBuiltinInfoFunc(b *testing.B) {
	benchmarkVectorizedBuiltinFunc(b, vecBuiltinInfoCases)
}
