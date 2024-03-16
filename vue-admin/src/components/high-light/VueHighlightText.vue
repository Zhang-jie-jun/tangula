<template>
  <div style="height: 480px; padding: 1px;">
    <span class="log-info-span" ref="showContent" v-highlight="{keyword, sensitive, overWriteStyle}">{{ data }}</span>
        <!-- <slot></slot> -->
    <!-- </span> -->
  </div>
</template>

<script>
import HighlightText from './HighlightText'
export default {
  props: {
    keyword: {
      type: [String, Array],
      required: true,
      default: 'ERROR'
    },
    sensitive: {
      type: Boolean,
      default: true
    },
    overWriteStyle: {
      type: Object,
      default: () => { return {} }
    },
    data: {
      type: String,
      required: true,
      default: ''
    }
  },
  watch: {
    data (newVal) {
      let element = this.$refs.showContent
      element.scrollTop = element.scrollHeight - element.clientHeight
    }
  },
  directives: {
    highlight: HighlightText
  }
}
</script>

<style lang="less">

.log-info-span {
  word-break:normal;
  // width:auto;
  width: 100%;
  height: 480px;
  display:block;
  white-space:pre-wrap;
  word-wrap: break-word;
  overflow: scroll;
  background: #353535;
  color: #FFFFFF;
  font-size: 14px;
}
</style>
