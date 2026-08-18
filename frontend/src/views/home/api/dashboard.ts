import dayjs from 'dayjs';
export interface ContentDataRecord {
  x: string;
  y: number;
}

export function queryContentData():any {
  const presetData = [58, 81, 53, 90, 64, 88];
  const getLineData = () => {
    const count = 6;
    return new Array(count).fill(0).map((el, idx) => ({
      x: dayjs()
        .day(idx - 2)
        .format('YYYY-MM-DD'),
      y: presetData[idx],
    }));
  };
  return successResponseWrap([...getLineData()]);
}
//组装返回数据
const successResponseWrap = (data: unknown) => {
  return {
    data,
    status: 'ok',
    msg: '请求成功',
    code: 20000,
  };
};
