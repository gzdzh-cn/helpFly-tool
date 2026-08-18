// 账单表格本地演示用的内存 mock 数据层（前端脚手架演示，不依赖后端）
export interface ConsumptionMock {
  key: string;
  taskId: string;
  tradeDate: string;
  channel: string;
  summary: string;
  incomeOrExpenseAmount: number;
  balance: number;
}

const channels = ['微信', '支付宝', '银行卡', '现金', '云闪付'];
const summaries = [
  '超市购物',
  '餐饮消费',
  '交通费',
  '话费充值',
  '工资入账',
  '水电费',
  '网购订单',
  '娱乐消费',
  '医疗支出',
  '理财收益',
];

function pad(n: number): string {
  return n < 10 ? `0${n}` : `${n}`;
}

function randomDate(): string {
  const y = 2025;
  const m = 1 + Math.floor(Math.random() * 12);
  const d = 1 + Math.floor(Math.random() * 28);
  const hh = pad(Math.floor(Math.random() * 24));
  const mm = pad(Math.floor(Math.random() * 60));
  const ss = pad(Math.floor(Math.random() * 60));
  return `${y}-${pad(m)}-${pad(d)} ${hh}:${mm}:${ss}`;
}

// 生成 n 条演示数据
export function genMockConsumptions(n = 86): ConsumptionMock[] {
  const list: ConsumptionMock[] = [];
  let balance = 10000;
  for (let i = 0; i < n; i += 1) {
    const isIncome = Math.random() > 0.7;
    const amount = Number((Math.random() * 2000 + 10).toFixed(2));
    balance = Number((balance + (isIncome ? amount : -amount)).toFixed(2));
    list.push({
      key: `M${Date.now()}_${i}`,
      taskId: 'demo-task',
      tradeDate: randomDate(),
      channel: channels[Math.floor(Math.random() * channels.length)],
      summary: summaries[Math.floor(Math.random() * summaries.length)],
      incomeOrExpenseAmount: isIncome ? amount : -amount,
      balance,
    });
  }
  return list;
}

export const channelOptions = channels;
