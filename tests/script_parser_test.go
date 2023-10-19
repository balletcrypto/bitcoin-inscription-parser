package parser

import (
	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/txscript"
	"testing"
)

func TestScriptWithInscription(t *testing.T) {
	t.Parallel()

	script, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		// If we push []byte{1} as content type tag, the script builder AddData method will convert this into small
		// number OP code: https://github.com/btcsuite/btcd/blob/master/txscript/scriptbuilder.go#L168
		// So we use two OP_DATA_1 with AddOp method instead
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddOp(txscript.OP_0).
		AddData([]byte("test script with inscription")).
		AddOp(txscript.OP_ENDIF).Script()

	scriptWithUnknownTag, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddData([]byte("unknown tag")). // [117 110 107 110 111 119 110 32 116 97 103]
		AddData([]byte("unknown data")).
		AddOp(txscript.OP_0).
		AddData([]byte("test script with inscription")).
		AddOp(txscript.OP_ENDIF).Script()

	scriptWithAdditionsAfterEndIf, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddOp(txscript.OP_0).
		AddData([]byte("test script with inscription")).
		AddOp(txscript.OP_ENDIF).
		AddOp(txscript.OP_CHECKSIG).
		Script()

	tests := []struct {
		testCase string
		script   []byte
		expected bool
	}{
		{
			testCase: "test script with inscription",
			script:   script,
			expected: true,
		},
		{
			testCase: "test script with unknown tag",
			script:   scriptWithUnknownTag,
			expected: true,
		},
		{
			testCase: "test script with additions after OP_ENDIF",
			script:   scriptWithAdditionsAfterEndIf,
			expected: true,
		},
	}

	for _, test := range tests {
		inscriptions := parser.ParseInscriptions(test.script)
		if len(inscriptions) != 0 {
			for i := range inscriptions {
				t.Logf("Find inscription with content type: %s, content length: %d", inscriptions[i].ContentType,
					inscriptions[i].ContentLength)
			}
		}
		exist := len(inscriptions) != 0
		if exist == test.expected {
			t.Logf("%s: test passed", test.testCase)
		} else {
			t.Errorf("%s: test failed", test.testCase)
		}
	}
}

func TestScriptWithMultipleInscriptions(t *testing.T) {
	script, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddOp(txscript.OP_0).
		AddData([]byte("test script with multiple inscriptions")).
		AddOp(txscript.OP_ENDIF).
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddOp(txscript.OP_0).
		AddData([]byte("test script with multiple inscriptions")).
		AddOp(txscript.OP_ENDIF).
		Script()

	test := struct {
		testCase string
		script   []byte
		expected bool
	}{
		testCase: "test script with multiple inscriptions",
		script:   script,
		expected: true,
	}

	inscriptions := parser.ParseInscriptions(test.script)
	if len(inscriptions) != 0 {
		for i := range inscriptions {
			t.Logf("Find inscription with content type: %s, content length: %d", inscriptions[i].ContentType,
				inscriptions[i].ContentLength)
		}
	}
	expected := len(inscriptions) == 2
	if expected == test.expected {
		t.Logf("%s: test passed", test.testCase)
	} else {
		t.Errorf("%s: test failed", test.testCase)
	}
}

func TestScriptWithIncompleteEnvelopHeader(t *testing.T) {
	t.Parallel()

	// No OP_FALSE header
	scriptWithNoOPFalse, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddOp(txscript.OP_0).
		AddData([]byte("test script with incomplete envelop header")).
		AddOp(txscript.OP_ENDIF).Script()

	// No OP_IF header
	scriptWithNoOPIf, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddOp(txscript.OP_0).
		AddData([]byte("test script with incomplete envelop header")).
		AddOp(txscript.OP_ENDIF).Script()

	// No ord header
	scriptWithNoOrd, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddOp(txscript.OP_0).
		AddData([]byte("test script with incomplete envelop header")).
		AddOp(txscript.OP_ENDIF).Script()

	tests := []struct {
		testCase string
		script   []byte
		expected bool
	}{
		{
			testCase: "test script with no OP_FALSE",
			script:   scriptWithNoOPFalse,
			expected: false,
		},
		{
			testCase: "test script with no OP_IF",
			script:   scriptWithNoOPIf,
			expected: false,
		},
		{
			testCase: "test script with no ord",
			script:   scriptWithNoOrd,
			expected: false,
		},
	}

	for _, test := range tests {
		inscriptions := parser.ParseInscriptions(test.script)
		if len(inscriptions) != 0 {
			for i := range inscriptions {
				t.Logf("Find inscription with content type: %s, content length: %d", inscriptions[i].ContentType,
					inscriptions[i].ContentLength)
			}
		}
		exist := len(inscriptions) != 0
		if exist == test.expected {
			t.Logf("%s: test passed", test.testCase)
		} else {
			t.Errorf("%s: test failed", test.testCase)
		}
	}
}

func TestScriptWithDuplicatedTag(t *testing.T) {
	t.Parallel()

	// Duplicated content type tag
	scriptWithDuplicatedContentTypeTag, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddOp(txscript.OP_0).
		AddData([]byte("test script with duplicated tag")).
		AddOp(txscript.OP_ENDIF).Script()

	// Unknown duplicated tag
	scriptWithUnknownDuplicatedTag, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddData([]byte("unknown tag")). // [117 110 107 110 111 119 110 32 116 97 103]
		AddData([]byte("unknown data")).
		AddData([]byte("unknown tag")).
		AddData([]byte("unknown data")).
		AddOp(txscript.OP_0).
		AddData([]byte("test script with duplicated tag")).
		AddOp(txscript.OP_ENDIF).Script()

	tests := []struct {
		testCase string
		script   []byte
		expected bool
	}{
		{
			testCase: "test script with duplicated content type tag",
			script:   scriptWithDuplicatedContentTypeTag,
			expected: false,
		},
		{
			testCase: "test script with unknown duplicated tag",
			script:   scriptWithUnknownDuplicatedTag,
			expected: false,
		},
	}

	for _, test := range tests {
		inscriptions := parser.ParseInscriptions(test.script)
		if len(inscriptions) != 0 {
			for i := range inscriptions {
				t.Logf("Find inscription with content type: %s, content length: %d", inscriptions[i].ContentType,
					inscriptions[i].ContentLength)
			}
		}
		exist := len(inscriptions) != 0
		if exist == test.expected {
			t.Logf("%s: test passed", test.testCase)
		} else {
			t.Errorf("%s: test failed", test.testCase)
		}
	}
}

func TestScriptWithOtherOpcodeBeforeEndIf(t *testing.T) {
	script, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddOp(txscript.OP_0).
		AddData([]byte("test script with other opcode before OP_ENDIF")).
		AddOp(txscript.OP_CHECKSIG).
		AddOp(txscript.OP_ENDIF).Script()

	test := struct {
		testCase string
		script   []byte
		expected bool
	}{
		testCase: "test script with other opcode before OP_ENDIF",
		script:   script,
		expected: false,
	}

	inscriptions := parser.ParseInscriptions(test.script)
	if len(inscriptions) != 0 {
		for i := range inscriptions {
			t.Logf("Find inscription with content type: %s, content length: %d", inscriptions[i].ContentType,
				inscriptions[i].ContentLength)
		}
	}
	exist := len(inscriptions) != 0
	if exist == test.expected {
		t.Logf("%s: test passed", test.testCase)
	} else {
		t.Errorf("%s: test failed", test.testCase)
	}
}

func TestScriptWithUnrecognizedEvenTag(t *testing.T) {
	script, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddData([]byte("test tag")). // [116 101 115 116 32 116 97 103]
		AddData([]byte("test data")).
		AddOp(txscript.OP_0).
		AddData([]byte("test script with unrecognized even tag")).
		AddOp(txscript.OP_ENDIF).Script()

	test := struct {
		testCase string
		script   []byte
		expected bool
	}{
		testCase: "test script with unrecognized even tag",
		script:   script,
		expected: true,
	}

	inscriptions := parser.ParseInscriptions(test.script)
	if len(inscriptions) != 0 {
		for i := range inscriptions {
			t.Logf("Find inscription with content type: %s, content length: %d", inscriptions[i].ContentType,
				inscriptions[i].ContentLength)
		}
	}
	exist := len(inscriptions) != 0
	if exist == test.expected {
		t.Logf("%s: test passed", test.testCase)
	} else {
		t.Errorf("%s: test failed", test.testCase)
	}
}

func TestScriptWithNoContentTypeTag(t *testing.T) {
	script, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_0).
		AddData([]byte("test script with no content type tag")).
		AddOp(txscript.OP_ENDIF).Script()

	test := struct {
		testCase string
		script   []byte
		expected bool
	}{
		testCase: "test script with no content type tag",
		script:   script,
		expected: true,
	}

	inscriptions := parser.ParseInscriptions(test.script)
	if len(inscriptions) != 0 {
		for i := range inscriptions {
			t.Logf("Find inscription with content type: %s, content length: %d", inscriptions[i].ContentType,
				inscriptions[i].ContentLength)
		}
	}
	exist := len(inscriptions) != 0
	if exist == test.expected {
		t.Logf("%s: test passed", test.testCase)
	} else {
		t.Errorf("%s: test failed", test.testCase)
	}
}

func TestScriptWithNoContentType(t *testing.T) {
	script, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_0).
		AddOp(txscript.OP_0).
		AddData([]byte("test script with no content type")).
		AddOp(txscript.OP_ENDIF).Script()

	test := struct {
		testCase string
		script   []byte
		expected bool
	}{
		testCase: "test script with no content type",
		script:   script,
		expected: true,
	}

	inscriptions := parser.ParseInscriptions(test.script)
	if len(inscriptions) != 0 {
		for i := range inscriptions {
			t.Logf("Find inscription with content type: %s, content length: %d", inscriptions[i].ContentType,
				inscriptions[i].ContentLength)
		}
	}
	exist := len(inscriptions) != 0
	if exist == test.expected {
		t.Logf("%s: test passed", test.testCase)
	} else {
		t.Errorf("%s: test failed", test.testCase)
	}
}

func TestScriptWithNoContentBody(t *testing.T) {
	script, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddOp(txscript.OP_ENDIF).Script()

	test := struct {
		testCase string
		script   []byte
		expected bool
	}{
		testCase: "test script with no content body",
		script:   script,
		expected: true,
	}

	inscriptions := parser.ParseInscriptions(test.script)
	if len(inscriptions) != 0 {
		for i := range inscriptions {
			t.Logf("Find inscription with content type: %s, content length: %d", inscriptions[i].ContentType,
				inscriptions[i].ContentLength)
		}
	}
	exist := len(inscriptions) != 0
	if exist == test.expected {
		t.Logf("%s: test passed", test.testCase)
	} else {
		t.Errorf("%s: test failed", test.testCase)
	}
}

func TestScriptWithInvalidDtaLength(t *testing.T) {
	script, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_66).
		AddOp(txscript.OP_0).
		AddData([]byte{0x62, 0x6f, 0x62, 0x2e, 0x73, 0x61, 0x74, 0x73, 0x0a}).
		AddOp(txscript.OP_ENDIF).Script()

	test := struct {
		testCase string
		script   []byte
		expected bool
	}{
		testCase: "test script with invalid data length",
		script:   script,
		expected: false,
	}

	inscriptions := parser.ParseInscriptions(test.script)
	if len(inscriptions) != 0 {
		for i := range inscriptions {
			t.Logf("Find inscription with content type: %s, content length: %d", inscriptions[i].ContentType,
				inscriptions[i].ContentLength)
		}
	}
	exist := len(inscriptions) != 0
	if exist == test.expected {
		t.Logf("%s: test passed", test.testCase)
	} else {
		t.Errorf("%s: test failed", test.testCase)
	}
}

func TestScriptWithZeroPush(t *testing.T) {
	script, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddOp(txscript.OP_0).
		AddOp(txscript.OP_ENDIF).Script()

	test := struct {
		testCase string
		script   []byte
		expected bool
	}{
		testCase: "test script with zero push",
		script:   script,
		expected: true,
	}

	inscriptions := parser.ParseInscriptions(test.script)
	if len(inscriptions) != 0 {
		for i := range inscriptions {
			t.Logf("Find inscription with content type: %s, content length: %d", inscriptions[i].ContentType,
				inscriptions[i].ContentLength)
		}
	}
	exist := len(inscriptions) != 0
	if exist == test.expected {
		t.Logf("%s: test passed", test.testCase)
	} else {
		t.Errorf("%s: test failed", test.testCase)
	}
}

func TestScriptWithMultiplePushes(t *testing.T) {
	script, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddOp(txscript.OP_0).
		AddData([]byte("test1 script with multiple pushes")).
		AddData([]byte("test2 script with multiple pushes")).
		AddOp(txscript.OP_ENDIF).Script()

	test := struct {
		testCase string
		script   []byte
		expected bool
	}{
		testCase: "test script with multiple pushes",
		script:   script,
		expected: true,
	}

	inscriptions := parser.ParseInscriptions(test.script)
	if len(inscriptions) != 0 {
		for i := range inscriptions {
			t.Logf("Find inscription with content type: %s, content length: %d", inscriptions[i].ContentType,
				inscriptions[i].ContentLength)
		}
	}
	exist := len(inscriptions) != 0
	if exist == test.expected {
		t.Logf("%s: test passed", test.testCase)
	} else {
		t.Errorf("%s: test failed", test.testCase)
	}
}

func TestScriptWithNoEndIf(t *testing.T) {
	t.Parallel()

	script1, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		AddOp(txscript.OP_0).
		AddData([]byte("test script with no END_IF")).
		Script()

	script2, _ := txscript.NewScriptBuilder().
		AddOps([]byte{txscript.OP_FALSE, txscript.OP_IF}).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte("text/plain;charset=utf-8")).
		Script()

	tests := []struct {
		testCase string
		script   []byte
		expected bool
	}{
		{
			testCase: "test script with no END_IF",
			script:   script1,
			expected: false,
		},
		{
			testCase: "test script with no END_IF",
			script:   script2,
			expected: false,
		},
	}

	for _, test := range tests {
		inscriptions := parser.ParseInscriptions(test.script)
		if len(inscriptions) != 0 {
			for i := range inscriptions {
				t.Logf("Find inscription with content type: %s, content length: %d", inscriptions[i].ContentType,
					inscriptions[i].ContentLength)
			}
		}
		exist := len(inscriptions) != 0
		if exist == test.expected {
			t.Logf("%s: test passed", test.testCase)
		} else {
			t.Errorf("%s: test failed", test.testCase)
		}
	}
}
