package validatortest

import (
	"fmt"
	"testing"

	"github.com/TheThingsIndustries/go-proto-validators/errors"
	"github.com/TheThingsIndustries/go-proto-validators/util"
)

var fullFieldMask = []string{"some_string", "embedded_mandatory", "embedded_not_mandatory", "field_mask", "embedded_mandatory.identifier", "embedded_mandatory.inner", "embedded_mandatory.field_mask", "embedded_mandatory.inner.id", "embedded_mandatory.inner.field_mask"}
var outerOnlyFieldMask = []string{"some_string", "field_mask"}
var outerAndMiddleOnlyFieldMask = []string{"some_string", "embedded_mandatory", "embedded_not_mandatory", "field_mask", "embedded_mandatory.identifier", "embedded_mandatory.field_mask"}
var fullFieldMaskWithEmbedded = []string{"some_string", "some_bytes", "embedded_repeated", "embedded_repeated_with_check", "embedded_repeated.identifier", "embedded_repeated.some_int", "embedded_repeated.field_mask", "embedded_repeated.inner", "embedded_repeated.inner.id", "embedded_repeated.inner.some_int", "embedded_repeated.inner.field_mask", "embedded_repeated_with_check.identifier", "embedded_repeated_with_check.some_int", "embedded_repeated_with_check.field_mask", "embedded_repeated_with_check.inner", "embedded_repeated_with_check.inner.id", "embedded_repeated_with_check.inner.some_int", "embedded_repeated_with_check.inner.field_mask"}
var outerOnlyFieldMaskWithEmbedded = []string{"some_string", "some_bytes", "embedded_repeated"}
var outerAndMiddleOnlyFieldMaskWithEmbedded = []string{"some_string", "some_bytes", "embedded_repeated", "embedded_repeated_with_check", "embedded_repeated.identifier", "embedded_repeated.some_int", "embedded_repeated.field_mask", "embedded_repeated_with_check.identifier", "embedded_repeated_with_check.some_int", "embedded_repeated_with_check.field_mask"}

func TestWithNilFieldMask(t *testing.T) {
	for _, tc := range []struct {
		Name                   string
		Message                interface{}
		FieldMask              []string
		ErrorExpected          bool
		ExpectedErrorFieldName string
		ExpectedErrorType      errors.Types
	}{
		{
			Name: "OuterWithValid",
			Message: &GenericTestMessage{
				SomeString: "outer",
				EmbeddedMandatory: &Embedded{
					Identifier: "middle",
					Inner: &InnerEmbedded{
						Id: "test-inner",
					},
				},
			},
			FieldMask:     nil,
			ErrorExpected: false,
		},
		{
			Name: "OuterWithInvalid",
			Message: &GenericTestMessage{
				SomeString: "&*^",
				EmbeddedMandatory: &Embedded{
					Identifier: "middle",
					Inner: &InnerEmbedded{
						Id: "test-inner",
					},
				},
			},
			FieldMask:     nil,
			ErrorExpected: false,
		},
		{
			Name: "MiddleWithInvalid",
			Message: &GenericTestMessage{
				SomeString: "outer",
				EmbeddedMandatory: &Embedded{
					Identifier: "&^&",
					Inner: &InnerEmbedded{
						Id: "test-inner",
					},
				},
			},
			FieldMask:     nil,
			ErrorExpected: false,
		},
		{
			Name: "InnerWithInvalid",
			Message: &GenericTestMessage{
				SomeString: "outer",
				EmbeddedMandatory: &Embedded{
					Identifier: "middle",
					Inner: &InnerEmbedded{
						Id: "",
					},
				},
			},
			FieldMask:     nil,
			ErrorExpected: false,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			var err error
			if msgVal, ok := tc.Message.(util.Validator); ok {
				err = msgVal.Validate(tc.FieldMask)
			} else {
				t.Fatal("Validator not found for message")
			}
			if tc.ErrorExpected && err != nil {
				if tc.ExpectedErrorFieldName != errors.GetFieldName(err.Error()) || tc.ExpectedErrorType.String() != errors.GetType(err.Error()) {
					t.Fatal(err)
				}
			} else if tc.ErrorExpected && err == nil {
				t.Fatal(fmt.Sprintf("Error %s Expected on field %s", tc.ExpectedErrorType, tc.ExpectedErrorFieldName))
			} else if !tc.ErrorExpected && err != nil {
				t.Fatal(err)
			} else {
				// Test Passed. Added for completeness
			}
		})
	}
}

func TestWithFullFieldMask(t *testing.T) {
	for _, tc := range []struct {
		Name                   string
		Message                interface{}
		FieldMask              []string
		ErrorExpected          bool
		ExpectedErrorFieldName string
		ExpectedErrorType      errors.Types
	}{
		{
			Name: "ValidMessage",
			Message: &GenericTestMessage{
				SomeString: "outer",
				EmbeddedMandatory: &Embedded{
					Identifier: "middle",
					Inner: &InnerEmbedded{
						Id: "inner",
					},
				},
			},
			FieldMask:     fullFieldMask,
			ErrorExpected: false,
		},
		{
			Name: "InnerInvalid",
			Message: &GenericTestMessage{
				SomeString: "outer",
				EmbeddedMandatory: &Embedded{
					Identifier: "middle",
					Inner: &InnerEmbedded{
						Id: "",
					},
				},
			},
			FieldMask:              fullFieldMask,
			ErrorExpected:          true,
			ExpectedErrorFieldName: "embedded_mandatory.inner.id",
			ExpectedErrorType:      errors.Types_STRING_NOT_EMPTY,
		},
		{
			Name: "MiddleInvalid",
			Message: &GenericTestMessage{
				SomeString: "outer",
				EmbeddedMandatory: &Embedded{
					Identifier: "&*^",
					Inner: &InnerEmbedded{
						Id: "",
					},
				},
			},
			FieldMask:              fullFieldMask,
			ErrorExpected:          true,
			ExpectedErrorFieldName: "embedded_mandatory.identifier",
			ExpectedErrorType:      errors.Types_STRING_REGEX,
		},
		{
			Name: "OuterInvalid",
			Message: &GenericTestMessage{
				SomeString: "&*sfsafsdfsfsdfdsfsdfsdf",
				EmbeddedMandatory: &Embedded{
					Identifier: "&HR",
					Inner: &InnerEmbedded{
						Id: "",
					},
				},
			},
			FieldMask:              fullFieldMask,
			ErrorExpected:          true,
			ExpectedErrorFieldName: "some_string",
			ExpectedErrorType:      errors.Types_STRING_REGEX,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			var err error
			if msgVal, ok := tc.Message.(util.Validator); ok {
				err = msgVal.Validate(tc.FieldMask)
			} else {
				t.Fatal("Validator not found for message")
			}
			if tc.ErrorExpected && err != nil {
				if tc.ExpectedErrorFieldName != errors.GetFieldName(err.Error()) || tc.ExpectedErrorType.String() != errors.GetType(err.Error()) {
					t.Fatal(err)
				}
			} else if tc.ErrorExpected && err == nil {
				t.Fatal(fmt.Sprintf("Error %s Expected on field %s", tc.ExpectedErrorType, tc.ExpectedErrorFieldName))
			} else if !tc.ErrorExpected && err != nil {
				t.Fatal(err)
			} else {
				// Test Passed. Added for completeness
			}
		})
	}
}

func TestWithPartialFieldMask(t *testing.T) {
	for _, tc := range []struct {
		Name                   string
		Message                interface{}
		FieldMask              []string
		ErrorExpected          bool
		ExpectedErrorFieldName string
		ExpectedErrorType      errors.Types
	}{
		{
			Name: "ValidOuterWithInvalidMiddleFMNotSet",
			Message: &GenericTestMessage{
				SomeString: "outer",
				EmbeddedMandatory: &Embedded{
					Identifier: "^&^",
					Inner: &InnerEmbedded{
						Id: "test-inner",
					},
				},
			},
			FieldMask:     outerOnlyFieldMask,
			ErrorExpected: false,
		},
		{
			Name: "ValidOuterWithInvalidMiddleFMNotSet1",
			Message: &GenericTestMessage{
				SomeString: "outer",
				EmbeddedMandatory: &Embedded{
					Identifier: "^&jfskfjdfjsdkfaf;sd^",
					Inner: &InnerEmbedded{
						Id: "test-inner",
					},
				},
			},
			FieldMask:     []string{"some_string", "field_mask", "embedded_mandatory", "embedded_mandatory.some_string"},
			ErrorExpected: false,
		},
		{
			Name: "ValidOuterWithInvalidMiddleFMSet",
			Message: &GenericTestMessage{
				SomeString: "outer",
				EmbeddedMandatory: &Embedded{
					Identifier: "^&^",
					Inner: &InnerEmbedded{
						Id: "test-inner",
					},
				},
			},
			FieldMask:              []string{"some_string", "field_mask", "embedded_mandatory", "embedded_mandatory.identifier"},
			ErrorExpected:          true,
			ExpectedErrorFieldName: "embedded_mandatory.identifier",
			ExpectedErrorType:      errors.Types_STRING_REGEX,
		},
		{
			Name: "ValidOuterWithOnlyOnefield",
			Message: &GenericTestMessage{
				SomeString:        "outer",
				EmbeddedMandatory: &Embedded{},
			},
			FieldMask:     []string{"some_string"},
			ErrorExpected: false,
		},
		{
			Name: "InvalidOuterWithOnlyOnefield",
			Message: &GenericTestMessage{
				SomeString:        "&%*$sd",
				EmbeddedMandatory: &Embedded{},
			},
			FieldMask:              []string{"some_string"},
			ErrorExpected:          true,
			ExpectedErrorFieldName: "some_string",
			ExpectedErrorType:      errors.Types_STRING_REGEX,
		},
		{
			Name: "InvalidOuterWithInvalidMiddle",
			Message: &GenericTestMessage{
				SomeString: "*&(&*#&$$",
				EmbeddedMandatory: &Embedded{
					Identifier: "^&^",
					Inner: &InnerEmbedded{
						Id: "test-inner",
					},
				},
			},
			FieldMask:              outerOnlyFieldMask,
			ErrorExpected:          true,
			ExpectedErrorFieldName: "some_string",
			ExpectedErrorType:      errors.Types_STRING_REGEX,
		},
		{
			Name: "ValidOuterAndMiddleWithInvalidInnerFMNotSet",
			Message: &GenericTestMessage{
				SomeString: "outer",
				EmbeddedMandatory: &Embedded{
					Identifier: "middle",
					Inner: &InnerEmbedded{
						Id: "test-inner",
					},
				},
			},
			FieldMask:     outerAndMiddleOnlyFieldMask,
			ErrorExpected: false,
		},
		{
			Name: "ValidOuterAndMiddleWithInvalidInnerFMNotSet1",
			Message: &GenericTestMessage{
				SomeString: "outer",
				EmbeddedMandatory: &Embedded{
					Identifier: "middle",
					Inner: &InnerEmbedded{
						Id: "test-inner",
					},
				},
			},
			FieldMask:     []string{"some_string", "embedded_mandatory", "embedded_not_mandatory", "field_mask", "embedded_mandatory.identifier", "embedded_mandatory.field_mask", "embedded_mandatory.inner", "embedded_mandatory.inner.id"},
			ErrorExpected: false,
		},
		{
			Name: "ValidOuterAndMiddleWithInvalidInnerFMSet",
			Message: &GenericTestMessage{
				SomeString: "outer",
				EmbeddedMandatory: &Embedded{
					Identifier: "middle",
					Inner: &InnerEmbedded{
						Id: "",
					},
				},
			},
			FieldMask:              []string{"some_string", "embedded_mandatory", "embedded_not_mandatory", "field_mask", "embedded_mandatory.identifier", "embedded_mandatory.field_mask", "embedded_mandatory.inner.id"},
			ErrorExpected:          true,
			ExpectedErrorFieldName: "embedded_mandatory.inner.id",
			ExpectedErrorType:      errors.Types_STRING_NOT_EMPTY,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			var err error
			if msgVal, ok := tc.Message.(util.Validator); ok {
				err = msgVal.Validate(tc.FieldMask)
			} else {
				t.Fatal("Validator not found for message")
			}
			if tc.ErrorExpected && err != nil {
				if tc.ExpectedErrorFieldName != errors.GetFieldName(err.Error()) || tc.ExpectedErrorType.String() != errors.GetType(err.Error()) {
					t.Fatal(err)
				}
			} else if tc.ErrorExpected && err == nil {
				t.Fatal(fmt.Sprintf("Error %s Expected on field %s", tc.ExpectedErrorType, tc.ExpectedErrorFieldName))
			} else if !tc.ErrorExpected && err != nil {
				t.Fatal(err)
			} else {
				// Test Passed. Added for completeness
			}
		})
	}
}

func TestWithRepeatedFields(t *testing.T) {
	for _, tc := range []struct {
		Name                   string
		Message                interface{}
		FieldMask              []string
		ErrorExpected          bool
		ExpectedErrorFieldName string
		ExpectedErrorType      errors.Types
	}{
		{
			Name: "ValidWithNilFieldMask",
			Message: &GenericTestMessageWithRepeated{
				SomeString:       "outer",
				EmbeddedRepeated: nil,
				EmbeddedRepeatedWithCheck: []*Embedded{
					&Embedded{
						Identifier: "middle",
						Inner: &InnerEmbedded{
							Id: "inner",
						},
					},
				},
			},
			FieldMask:     []string{},
			ErrorExpected: false,
		},
		{
			Name: "InvalidWithNilFieldMask",
			Message: &GenericTestMessageWithRepeated{
				SomeString:       "outer",
				EmbeddedRepeated: nil,
				EmbeddedRepeatedWithCheck: []*Embedded{
					&Embedded{
						Identifier: "middle",
						Inner: &InnerEmbedded{
							Id: "",
						},
					},
				},
			},
			FieldMask:     nil,
			ErrorExpected: false,
		},
		{
			Name: "NoFieldWithNilFieldMask",
			Message: &GenericTestMessageWithRepeated{
				SomeString:                "outer",
				EmbeddedRepeated:          nil,
				EmbeddedRepeatedWithCheck: nil,
			},
			FieldMask:              []string{},
			ErrorExpected:          true,
			ExpectedErrorFieldName: "embedded_repeated_with_check",
			ExpectedErrorType:      errors.Types_REPEATED_COUNT_MIN,
		},
		{
			Name: "ValidWithFullFieldMask",
			Message: &GenericTestMessageWithRepeated{
				SomeString: "outer",
				EmbeddedRepeated: []*Embedded{
					&Embedded{
						Identifier: "middle",
						Inner: &InnerEmbedded{
							Id: "inner",
						},
					},
				},
				EmbeddedRepeatedWithCheck: []*Embedded{
					&Embedded{
						Identifier: "middle",
						Inner: &InnerEmbedded{
							Id: "inner",
						},
					},
				},
			},
			FieldMask:     fullFieldMaskWithEmbedded,
			ErrorExpected: false,
		},
		{
			Name: "InvalidWithFullFieldMask",
			Message: &GenericTestMessageWithRepeated{
				SomeString: "outer",
				EmbeddedRepeated: []*Embedded{
					&Embedded{
						Identifier: "middle",
						Inner: &InnerEmbedded{
							Id: "",
						},
					},
				},
				EmbeddedRepeatedWithCheck: []*Embedded{
					&Embedded{
						Identifier: "middle",
						Inner: &InnerEmbedded{
							Id: "inner",
						},
					},
				},
			},
			FieldMask:              fullFieldMaskWithEmbedded,
			ErrorExpected:          true,
			ExpectedErrorFieldName: "embedded_repeated.inner.id",
			ExpectedErrorType:      errors.Types_STRING_NOT_EMPTY,
		},
		{
			Name: "InvalidRepeatedWithFullFieldMask",
			Message: &GenericTestMessageWithRepeated{
				SomeString: "outer",
				EmbeddedRepeated: []*Embedded{
					&Embedded{
						Identifier: "middle",
						Inner: &InnerEmbedded{
							Id: "inner",
						},
					},
					&Embedded{
						Identifier: "middle",
						Inner: &InnerEmbedded{
							Id: "inner",
						},
					},
				},

				EmbeddedRepeatedWithCheck: []*Embedded{
					&Embedded{
						Identifier: "middle",
						Inner: &InnerEmbedded{
							Id: "",
						},
					},
				},
			},
			FieldMask:              fullFieldMaskWithEmbedded,
			ErrorExpected:          true,
			ExpectedErrorFieldName: "embedded_repeated_with_check.inner.id",
			ExpectedErrorType:      errors.Types_STRING_NOT_EMPTY,
		},
		{
			Name: "InvalidWithPartialFieldMask",
			Message: &GenericTestMessageWithRepeated{
				SomeString: "outer",
				EmbeddedRepeated: []*Embedded{
					&Embedded{
						Identifier: "middle",
						Inner: &InnerEmbedded{
							Id: "inner",
						},
					},
				},
				EmbeddedRepeatedWithCheck: []*Embedded{
					&Embedded{
						Identifier: "middle",
						Inner: &InnerEmbedded{
							Id: "inner",
						},
					},
				},
			},
			FieldMask:     outerOnlyFieldMaskWithEmbedded,
			ErrorExpected: false,
		},
		{
			Name: "InvalidWithPartialFieldMask2",
			Message: &GenericTestMessageWithRepeated{
				SomeString: "outer",
				EmbeddedRepeated: []*Embedded{
					&Embedded{
						Identifier: "middle",
						Inner: &InnerEmbedded{
							Id: "inner",
						},
					},
				},
				EmbeddedRepeatedWithCheck: []*Embedded{
					&Embedded{
						Identifier: "middle",
						Inner: &InnerEmbedded{
							Id: "inner",
						},
					},
				},
			},
			FieldMask:     outerAndMiddleOnlyFieldMaskWithEmbedded,
			ErrorExpected: false,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			var err error
			if msgVal, ok := tc.Message.(util.Validator); ok {
				err = msgVal.Validate(tc.FieldMask)
			} else {
				t.Fatal("Validator not found for message")
			}
			if tc.ErrorExpected && err != nil {
				if tc.ExpectedErrorFieldName != errors.GetFieldName(err.Error()) || tc.ExpectedErrorType.String() != errors.GetType(err.Error()) {
					t.Fatal(err)
				}
			} else if tc.ErrorExpected && err == nil {
				t.Fatal(fmt.Sprintf("Error %s Expected on field %s", tc.ExpectedErrorType, tc.ExpectedErrorFieldName))
			} else if !tc.ErrorExpected && err != nil {
				t.Fatal(err)
			} else {
				// Test Passed. Added for completeness
			}
		})
	}
}

// func TestSomething(t *testing.T) {
// 	msg := GenericTestMessage{
// 		SomeString: "hello",
// 		EmbeddedMandatory: &Embedded{
// 			Identifier: "&*&",
// 		},
// 	}

// 	fmt.Println(msg.Validate([]string{"some_string", "embedded_mandatory.identifier"}))
// }
