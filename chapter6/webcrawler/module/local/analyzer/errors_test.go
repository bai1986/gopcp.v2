package analyzer

import (
	"testing"

	"gopcp.v2/chapter6/webcrawler/errors"
)

func TestErrorGenError(t *testing.T) {
	simpleErrMsg := "testing error"
	expectedErrType := errors.ERROR_TYPE_ANALYZER
	err := genError(simpleErrMsg)
	ce, ok := err.(errors.CrawlerError)
	if !ok {
		t.Fatalf("Inconsistent error type: expected: %T, actual: %T",
			errors.NewCrawlerError("", ""), err)
	}
	if ce.Type() != expectedErrType {
		t.Fatalf("Inconsistent error type string: expected: %q, actual: %q",
			expectedErrType, ce.Type())
	}
	expectedErrMsg := "crawler error: analyzer error: " + simpleErrMsg
	if ce.Error() != expectedErrMsg {
		t.Fatalf("Inconsistent error message: expected: %q, actual: %q",
			expectedErrMsg, ce.Error())
	}
}

func TestErrorGenErrorr(t *testing.T) {
	simpleErrMsg := "testing error"
	expectedErrType := errors.ERROR_TYPE_ANALYZERR
	err := genErrorr(simpleErrMsg)
	ce, ok := err.(errors.CrawlerErrorr)
	if !ok {
		t.Fatalf("inconsistent error type expected : %T, actual :%T",
			errors.NewCrawlerErrorr("",""),err)
	}
	if ce.Type() != expectedErrType {
		t.Fatalf("inconsistent error type string: expected: %q,actual:%q",expectedErrType,ce.Type())
	}
	expectedErrMsg := "crawler error:analyzer error:"+simpleErrMsg
	if ce.Error() != expectedErrMsg {
		t.Fatalf("inconsistent error message :expected: %q,actual: %q",expectedErrMsg,ce.Error())
	}
}
