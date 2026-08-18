import { DbService } from '/#/helpfly/internal/service';
import type { TaskPayload, TaskRecord } from '/#/helpfly/internal/service/models';

export type TaskFormPayload = Omit<TaskPayload, 'openingBalance' | 'transactionCount'> & {
  openingBalance?: number | null;
  transactionCount?: number | null;
};

const normalizePayload = (payload: TaskFormPayload): TaskPayload => ({
  ...(payload as any),
  openingBalance: payload.openingBalance ?? null,
  transactionCount: payload.transactionCount ?? null,
});

export function createTask(payload: TaskFormPayload) {
  return DbService.CreateTask(normalizePayload(payload));
}

export function updateTask(record: TaskRecord) {
  return DbService.UpdateTask(record);
}

export function deleteTasks(keys: string[]) {
  return DbService.DeleteTasks(keys);
}

export function listTasks() {
  return DbService.ListTasks();
}

export function searchTasks(keyword: string) {
  return DbService.SearchTasks(keyword);
}

export function testPromiose(key: string) {
  return DbService.TestPromiose(key);
}
