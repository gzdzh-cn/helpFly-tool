package service

import (
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func TestTaskValidationKeepsChineseMessages(t *testing.T) {
	err := validateTaskPayload(TaskPayload{})
	if err == nil || err.Error() != "任务名称不能为空" {
		t.Fatalf("unexpected task validation error: %v", err)
	}
	if gerror.Code(err) != gcode.CodeInvalidParameter {
		t.Fatalf("expected invalid parameter code, got %v", gerror.Code(err))
	}
}

func TestParameterOptionValidationKeepsChineseMessages(t *testing.T) {
	err := validateOptionList(nil, "银行选项不能为空")
	if err == nil || err.Error() != "银行选项不能为空" {
		t.Fatalf("unexpected option validation error: %v", err)
	}
	if gerror.Code(err) != gcode.CodeInvalidParameter {
		t.Fatalf("expected invalid parameter code, got %v", gerror.Code(err))
	}
}
