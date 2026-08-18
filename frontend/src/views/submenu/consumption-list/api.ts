// 账单表格数据访问层：根据 dataMode 在「后端」与「本地 mock」之间切换
import { DbService } from '/#/helpfly/internal/service';
import { genMockConsumptions, type ConsumptionMock } from './mock';

export type DataMode = 'backend' | 'mock';

// 内存 mock 存储（仅在 mock 模式下使用）
let mockStore: ConsumptionMock[] = [];

function toRecord(row: ConsumptionMock) {
  return {
    key: row.key,
    taskId: row.taskId,
    tradeDate: row.tradeDate,
    channel: row.channel,
    summary: row.summary,
    incomeOrExpenseAmount: row.incomeOrExpenseAmount,
    balance: row.balance,
  };
}

function fromRecord(r: any): ConsumptionMock {
  return {
    key: r.key || '',
    taskId: r.taskId || '',
    tradeDate: r.tradeDate || '',
    channel: r.channel || '',
    summary: r.summary || '',
    incomeOrExpenseAmount: Number(r.incomeOrExpenseAmount || 0),
    balance: Number(r.balance || 0),
  };
}

export async function fetchList(mode: DataMode, keyword = ''): Promise<ConsumptionMock[]> {
  if (mode === 'mock') {
    if (mockStore.length === 0) mockStore = genMockConsumptions();
    const list = keyword ? mockStore.filter((i) => JSON.stringify(i).includes(keyword)) : mockStore;
    return [...list];
  }
  const res = keyword
    ? await DbService.SearchConsumptions('', keyword)
    : await DbService.ListConsumptions('');
  return (res || []).map(fromRecord);
}

export async function addRecord(mode: DataMode, row: ConsumptionMock): Promise<ConsumptionMock> {
  if (mode === 'mock') {
    const created = { ...row, key: `M${Date.now()}_new` };
    mockStore.unshift(created);
    return created;
  }
  const res = await DbService.AddConsumption(toRecord(row));
  return fromRecord(res);
}

export async function updateRecord(mode: DataMode, row: ConsumptionMock): Promise<ConsumptionMock> {
  if (mode === 'mock') {
    const idx = mockStore.findIndex((i) => i.key === row.key);
    if (idx >= 0) mockStore[idx] = { ...row };
    return { ...row };
  }
  const res = await DbService.UpdateConsumption(toRecord(row));
  return fromRecord(res);
}

export async function deleteRecords(mode: DataMode, keys: string[]): Promise<void> {
  if (mode === 'mock') {
    mockStore = mockStore.filter((i) => !keys.includes(i.key));
    return;
  }
  await DbService.DeleteConsumptions(keys);
}

// 高级搜索：本地二次过滤（后端 SearchConsumptions 仅支持关键字）
export function applyAdvancedFilter(
  list: ConsumptionMock[],
  cond: { keyword?: string; channel?: string; type?: string; dateRange?: string[]; min?: number; max?: number },
): ConsumptionMock[] {
  return list.filter((i) => {
    if (cond.keyword && !JSON.stringify(i).includes(cond.keyword)) return false;
    if (cond.channel && i.channel !== cond.channel) return false;
    if (cond.type === 'income' && i.incomeOrExpenseAmount <= 0) return false;
    if (cond.type === 'expense' && i.incomeOrExpenseAmount >= 0) return false;
    if (cond.dateRange && cond.dateRange.length === 2) {
      const d = i.tradeDate.slice(0, 10);
      if (d < cond.dateRange[0] || d > cond.dateRange[1]) return false;
    }
    if (cond.min != null && Math.abs(i.incomeOrExpenseAmount) < cond.min) return false;
    if (cond.max != null && Math.abs(i.incomeOrExpenseAmount) > cond.max) return false;
    return true;
  });
}
