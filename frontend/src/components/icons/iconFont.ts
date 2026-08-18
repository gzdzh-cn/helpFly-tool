import { h } from 'vue'
import { Icon } from '@arco-design/web-vue';

const IconFont = Icon.addFromIconFontCn({ 
  src: 'https://at.alicdn.com/t/c/font_5075063_6uzohei6hip.js',
  extraProps:{
    type: 'icon-kuangjia',
    style: {
      fontSize: '18px',
    },
  }
});
const DynamicIconFont = (props:any) => {
  return h(IconFont, { type: props.type || 'icon-kuangjia' ,style:{fontSize:props.size+"px"|| '18px'}})
}

export default DynamicIconFont
