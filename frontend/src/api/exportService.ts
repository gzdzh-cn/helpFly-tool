import { DbService } from '/#/helpfly/internal/service';

export interface ExportRecord {
  key: string;
  taskId: string;
  filePath: string;
  createdAt: string;
}

export function listExports() {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.ListExports();
}

export function searchExports(keyword: string) {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.SearchExports(keyword);
}

export function deleteExports(keys: string[]) {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.DeleteExports(keys);
}

export function getDefaultExportPath() {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.GetDefaultExportPath();
}

export function selectExportPath(defaultPath?: string) {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.SelectExportPath(defaultPath || '');
}

export function openFolder(path: string) {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.OpenFolder(path);
}

export function getExportPath() {
  // @ts-ignore - 方法将在重新生成绑定后可用
  return DbService.GetExportPath();
}
