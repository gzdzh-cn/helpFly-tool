import { defHttp } from '@/utils/http';
import { FilterParam } from '/@/utils';
enum Api {
    // getList = '/createcode/product/getList',
    getData= '/wails/httprequest/getData',
    postData = '/wails/httprequest/postData',
    // upStatus = '/createcode/product/upStatus',
    // del = '/createcode/product/del',
}

//列表数据
// export function getList(params: any) {
//   params=FilterParam(params)
//   return defHttp.get({ url: Api.getList, params:params }, { errorMessageMode: 'message' });
// }
//获取内容
export function getData(params: object) {
  return defHttp.get({ url: Api.getData, params:params }, { errorMessageMode: 'message' });
}
//提交数据
export function postData(params: object) {
    return defHttp.post({ url: Api.postData, params:params}, { errorMessageMode: 'message' });
}
// //更新状态
// export function upStatus(params: object) {
//     return defHttp.post({ url: Api.upStatus, params:params}, { errorMessageMode: 'message' });
// }
// //删除数据
// export function del(params: object) {
//     return defHttp.delete({ url: Api.del, params:params}, { errorMessageMode: 'message' });
// }
/**数据类型 */
export interface DataItem {
    id:number,
    name: string;
}
