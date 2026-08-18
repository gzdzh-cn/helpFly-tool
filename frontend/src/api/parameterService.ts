import { DbService } from '/#/helpfly/internal/service';
import type { SystemParameters as SystemParametersModel } from '/#/helpfly/internal/service/models';

// 后端返回的参数（数组格式）
export interface SystemParameters {
  banks: string[];
  depositTypes: string[];
  cashExchanges: string[];
  currencies: string[];
  channels: string[];
  summaries: string[];
  regions: string[];
  exportPath: string;
  addWatermark: boolean;
  watermarkPath: string;
}

// 前端表单使用的参数（字符串格式，逗号分隔）
export interface SystemParametersForm {
  banks: string;
  depositTypes: string;
  cashExchanges: string;
  currencies: string;
  channels: string;
  summaries: string;
  regions: string;
  exportPath: string;
  addWatermark: boolean;
  watermarkPath: string;
}

export function getSystemParameters(): Promise<SystemParameters> {
  return DbService.GetSystemParameters().then((result: SystemParametersModel) => {
    return {
      banks: result.banks || [],
      depositTypes: result.depositTypes || [],
      cashExchanges: result.cashExchanges || [],
      currencies: result.currencies || [],
      channels: result.channels || [],
      summaries: result.summaries || [],
      regions: result.regions || [],
      exportPath: result.exportPath || '',
      addWatermark: result.addWatermark || false,
      watermarkPath: result.watermarkPath || '',
    };
  });
}

export function saveSystemParameters(params: SystemParameters): Promise<string> {
  // 直接传递对象，Wails 会自动处理类型转换
  return DbService.SaveSystemParameters(params as any);
}

export function selectWatermarkImage(): Promise<string> {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.SelectWatermarkImage();
}
