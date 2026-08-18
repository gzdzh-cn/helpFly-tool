// 表单演示用的本地 mock 数据
export interface DemoForm {
  id: string;
  name: string;
  category: string;
  email: string;
  phone: string;
  gender: string;
  birthday: string;
  amount: number;
  score: number;
  switchOn: boolean;
  tags: string[];
  remark: string;
}

export const categoryOptions = ['产品', '研发', '设计', '运营', '市场'];
export const tagOptions = ['紧急', '重要', '待跟进', '已完成', 'VIP'];

export function emptyForm(): DemoForm {
  return {
    id: '',
    name: '',
    category: '',
    email: '',
    phone: '',
    gender: '',
    birthday: '',
    amount: 0,
    score: 3,
    switchOn: true,
    tags: [],
    remark: '',
  };
}
