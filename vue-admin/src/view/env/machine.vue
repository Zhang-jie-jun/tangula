<template>
  <div class="layout">
    <Layout>
      <Content>
        <Menu mode="horizontal" :active-name="subItemIndex">
          <MenuItem :name="index+1" v-for="(item, index) in subMenuList" :key="item.id">
            <div v-on:click="changeSubItem(index+1, $event)">
              <Icon :type="item.icon"></Icon>
              {{item.name}}
            </div>
          </MenuItem>
        </Menu>
        <div style="min-height: 200px;">
          <div v-if="subItemIndex === 1">
            <platformPrivate v-if="isRouterAlive" />
          </div>
          <div v-if="subItemIndex === 2">
            <platformPublic v-if="isRouterAlive"/>

          </div>
        </div>
      </Content>
    </Layout>
    <Spin fix v-if="showSpin">
      <Icon type="ios-loading" size=18 class="demo-spin-icon-load" ></Icon>
      <div>Loading</div>
    </Spin>
  </div>
</template>

<script>
// @ is an alias to /src
import { mapState } from 'vuex'
import platformPublic from '@/components/machine/machinePublic.vue'
import platformPrivate from '@/components/machine/machinePrivate.vue'

export default {
  name: 'machine',
  components: {
    platformPublic,
    platformPrivate
  },
  data () {
    return {
      isRouterAlive: true,
      tvData: [], // 目标版本
      showSpin: false,
      menuIndex: 1,
      subItemIndex: 1,
      menuList: [],
      subMenuList: [
        { id: 1, name: '私有主机', icon: 'md-desktop' },
        { id: 2, name: '公共主机', icon: 'md-desktop' }
      ],
      model1: ''
    }
  },
  computed: {
    ...mapState({
      permits: state => state.user.permits,
      defBusi: state => state.user.userBusi,
      userName: state => state.user.userName
    })
  },
  methods: {
    changeSubItem: function (subItemIndex) {
      this.subItemIndex = subItemIndex
    },
    changeMenuItem: function (menuIndex, target) {
      this.subItemIndex = 1
      this.menuIndex = menuIndex
      this.selectObj.target = target
      this.reload()
    },
    // 刷新组件
    reload () {
      this.isRouterAlive = false
      this.$nextTick(() => (this.isRouterAlive = true))
    }
  },
  mounted: function () {
  }
}
</script>
<style scoped>
.layout{
  border: 1px solid #d7dde4;
  background: #f5f7f9;
  position: relative;
  border-radius: 4px;
  overflow: hidden;
}
.layout-Select{
  width: 200px;
  float: left;
  position: relative;
  top: 15px;
  left: 20px;
}

.layout-footer-center{
  text-align: center;
}
.ivu-menu{
  z-index: initial;
}
.menu-num{

}
/**
.layout-nav .ivu-menu-item-active{
    color: red !important;
}
*/
.layout-header ul{
  height: 100%;
}
.layout-header .ivu-menu{
  display: flex;

}
.layout-select{
  width: 200px;
}
.layout-nav{
  flex: 1;
  margin-left: 50px;
  display: flex;
  overflow-x: auto;
}
.layout-nav::-webkit-scrollbar {
  display: none;/*隐藏滚动条*/
}
.layout-nav li{
  position: relative;
  height: 100%;
  flex-shrink: 0;
  white-space: nowrap;
}
.menu-li-item{
  height: 100%;
}

.menu-li-item span{
  display: inline-block;
  position: relative;
  height: 100%;
}
.menu-li-item span i{
  position: absolute;
  top: 10px;
  right: -20px;
  display: inline-block;
  height: 20px;
  line-height: 20px;
  font-size: 12px;
  background-color: #2d8cf0;
  border-radius: 100%;
  width: 20px;
  color: #fff;
  text-align: center;
  font-style: normal;
  font-weight:normal;
}
.layout-header .ivu-menu-item-active{
  font-weight: bold;
}
.ivu-layout-header{
  padding: 0;
}
</style>
