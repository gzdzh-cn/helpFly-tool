import { DbService } from '/#/helpfly/internal/service';
import type { TaskRecord } from '/#/helpfly/internal/service/models';

export interface ConsumptionRecord {
  key: string;
  taskId: string;
  tradeDate: string;
  account: string;
  storageType: string;
  serialNumber: number;
  currency: string;
  cashOrRemit: string;
  summary: string;
  region: string;
  incomeOrExpenseAmount: number;
  balance: number;
  channel: string;
}

export interface TaskProgress {
  taskId: string;
  current: number;
  total: number;
  percent: number;
  isRunning: boolean;
  isCanceled: boolean;
  error?: string;
}

export function generateConsumptions(task: TaskRecord, doType: string = ''): Promise<string> {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.GenerateConsumptions(task, doType);
}

export function listConsumptions(taskId?: string) {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.ListConsumptions(taskId || '');
}

export function searchConsumptions(taskId: string, keyword: string) {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.SearchConsumptions(taskId, keyword);
}

export function deleteConsumptions(keys: string[]) {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.DeleteConsumptions(keys);
}

export function deleteConsumptionsByTask(taskId: string) {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.DeleteConsumptionsByTaskID(taskId);
}

export function getTaskProgress(taskId: string) {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.GetTaskProgress(taskId);
}

export function cancelTask(taskId: string) {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.CancelTask(taskId);
}

export function getDeleteProgress(taskId: string) {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.GetDeleteProgress(taskId);
}

export function exportConsumptionsToPDF(taskId: string, keys: string[]): Promise<string> {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.ExportConsumptionsToPDF(taskId, keys);
}

export function exportConsumptionsToExcel(taskId: string, keys: string[]): Promise<string> {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.ExportConsumptionsToExcel(taskId, keys);
}
