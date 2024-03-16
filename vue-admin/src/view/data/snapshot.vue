<template>
  <div>
    <Modal
      :title="this.backupName +'-快照列表'"
      v-model="showModal"
      :mask-closable="false"
      :closable="false"
      width="1400"
      @on-ok="closeModal"
      @on-cancel="closeModal"
    >
      <!-- 主体部分 -->
      <div>
        <Row>
          <Row class="search-margin-botton">
            <Col span="4">
              <Button type="info"  @click.prevent="openCreate">创建快照</Button>
            </Col>
            <Col span="12">
              <Row>
                <Col span="6" style="margin-right: 10px;">
                  <Input placeholder="名称" v-model="searchObj.snapshotName"/>
                </Col>
                <Button class="right-but" type="primary" @click="search">搜索</Button>
              </Row>
            </Col>
          </Row>
          <Table searchable border :columns="columnsTitle" :data="orders" stripe size="small">
            <template slot-scope="{ row }" slot="action">
              <div>
                <Button
                  type="info"
                  size="small"
                  style="margin-right: 5px"
                  @click="openBackup(row)"
                >生成副本</Button>
                <Button
                  type="info"
                  size="small"
                  style="margin-right: 5px"
                  @click="openMirror(row)"
                >生成镜像</Button>
                <Poptip
                  confirm
                  @on-ok="rollback(row)"
                  title="确定要回滚吗？">
                  <Button
                    type="warning"
                    size="small"
                    style="margin-right: 5px"
                  >回滚
                  </Button>
                </Poptip>
                <Poptip
                  confirm
                  @on-ok="destroy(row)"
                  title="确定要删除快照吗？">
                  <Button
                    type="error"
                    size="small"
                  >删除
                  </Button>
                </Poptip>
              </div>
            </template>
          </Table>
          <div class="table-page">
            <div class="table-page-position">
              <Page :total="total" :current="currentPage" show-total show-elevator @on-change="changePage"></Page>
            </div>
          </div>
          <Spin fix v-if="showSpin">
            <Icon type="ios-loading" size=18 class="demo-spin-icon-load" ></Icon>
            <div>Loading</div>
          </Spin>
        </Row>
      </div>
      <div slot="footer">
        <Button  @click="closeModal">关闭</Button>
      </div>
    </Modal>
    <Modal
      v-model="crtModal"
      draggable
      title="创建快照"
      @on-cancel="cancel">
      <Form :rules="ruleValidate" ref="formCrtPlat" :model="formCrtPlat">
        <FormItem label="名称" prop="name">
          <Input v-model="formCrtPlat.name" placeholder="请输入名称"/>
        </FormItem>
        <FormItem label="描述" prop="desc">
          <Input v-model="formCrtPlat.desc" placeholder="请输描述"/>
        </FormItem>
      </Form>
      <div slot="footer">
        <Button type="primary" :disabled="isdisabledFn" v-preventReClick @click="crtPlat('formCrtPlat')">确定</Button>
      </div>
    </Modal>
    <Modal
      v-model="mirrorModal"
      draggable
      title="生成镜像"
      @on-cancel="cancel">
      <Form :rules="ruleMirrorValidate" ref="formMirror" :model="formMirror">
        <FormItem label="名称" prop="name">
          <Input v-model="formMirror.name" placeholder="请输入名称"/>
        </FormItem>
        <FormItem label="描述" prop="desc">
          <Input v-model="formMirror.desc" placeholder="请输入描述"/>
        </FormItem>
      </Form>
      <div slot="footer">
        <Button type="primary" :disabled="isdisabledFn" v-preventReClick @click="crtMirror('formMirror')">确定</Button>
      </div>
    </Modal>
    <Modal
      draggable
      v-model="backupModal"
      title="生成副本"
      @on-cancel="cancel">
      <Form :rules="ruleMirrorValidate" ref="formBackup" :model="formBackup">
        <FormItem label="名称" prop="name">
          <Input v-model="formBackup.name" placeholder="请输入名称"/>
        </FormItem>
        <FormItem label="描述" prop="desc">
          <Input v-model="formBackup.desc" placeholder="请输入描述"/>
        </FormItem>
      </Form>
      <div slot="footer">
        <Button type="primary" :disabled="isdisabledFn" v-preventReClick @click="crtBackup('formBackup')">确定</Button>
      </div>
    </Modal>
  </div>
</template>
<script>
import { mapState } from 'vuex'
import {
  crtSnapshot,
  deleteBackup,
  deleteSnapshot,
  getSnapshotList, rollback,
  snapshotToBackup,
  snapshotToMirror
} from "@/api/backup";
import {httpRequest} from "@/libs/tools";

export default {
  name: 'snapshot',
  // 子组件自己的数据定义，可读可写
  props: {
    show: Boolean,
    modify: {
      default: false
    },
    backupName: String,
    backupId: Number
  },
  data () {
    return {
      isdisabledFn:false,
      // 搜索
      searchObj: {
        snapshotName: ''
      },
      formCrtPlat: {
        name: '',
        desc: ''
      },
      formMirror: {
        snapId: '',
        name: '',
        desc: ''
      },
      formBackup: {
        snapId: '',
        name: '',
        desc: ''
      },
      crtModal: false,
      showSpin: false,
      ruleValidate: {
        // reviewConclusion: [{ required: true, message: '请输入', trigger: 'change' }],
        // lastDate: [{ required: true, message: '请输入', trigger: 'change', type: 'date' }],
        name: [{ required: true, message: '请输入', trigger: 'blur' }]
      },
      ruleMirrorValidate: {
        name: [{ required: true, message: '请输入', trigger: 'blur' }]
      },
      searchChanged: false,
      mirrorModal: false,
      backupModal: false,
      // 需求详情模态框
      total: 0,
      currentPage: 1,
      orders: [],
      showModal: false,
      columnsTitle: [
        {
          title: 'ID',
          width: 70,
          key: 'id',
          align: 'center'
        },
        {
          title: '名称',
          key: 'name',
          align: 'center'
        },
        {
          title: '描述',
          key: 'desc',
          align: 'center'
        },
        {
          title: 'uuid',
          key: 'uuid',
          align: 'center'
        },
        {
          title: '创建时间',
          key: 'created_time',
          align: 'center'
        },

        {
          title: '操作',
          key: '',
          align: 'left',
          slot: 'action'
        }
      ]
    }
  },
  computed: {
    ...mapState({
      access: state => state.user.access,
      userName: state => state.user.userName
    }),
    showSelf () {
      return this.show
    }
  },
  watch: {
    showSelf (newVal, oldVal) {
      if(newVal){
        this.getData()
      }
      this.showModal = newVal

    }
  },
  mounted () {
    //this.getData()
  },
  methods: {
    openCreate () {
      this.crtModal = true
    },
    getData (page = 1) {
      this.showSpin = true
      const params = Object.assign({}, {
        id: Number(this.backupId),
        filter: this.searchObj.snapshotName,
      })

      httpRequest(this, getSnapshotList, [page, params]).then(res => {
        this.showSpin = false
        if (res.code !== 200) {
          this.$Notice.error({
            title: '错误',
            desc: res.message
          })
          return
        }
        if(res.response.data){
          this.orders = res.response.data
          this.total = res.response.totalNum
        }else {
          this.orders = []
          this.total = 0
        }
        this.currentPage = page

      })
    },
    crtPlat (name) {
      this.$refs[name].validate((valid) => {
        if (valid) {
          this.isdisabledFn=true
          const param = {
            id: Number(this.backupId),
            name: this.formCrtPlat.name,
            desc: this.formCrtPlat.desc
          }
          this.$Notice.info({
            title: '请求中……',
            desc: '正在创建快照，请等待',
            duration:10
          })
          crtSnapshot(param).then(({ data: { code, message } }) => {
            if (code !== 200) {
              this.$Notice.error({
                title: '错误',
                desc: message
              })
              this.isdisabledFn=false
              return
            }
            this.$Notice.success({
              title: '成功',
              desc: '操作成功'
            })
            this.crtModal = false
            this.isdisabledFn=false
            this.getData()
          }).catch(err => {
            if (err && err.response) {
              // console.error(err)
              this.$Notice.error({
                title: '错误',
                desc: err.response.data.msg
              })
              this.isdisabledFn=false
            }
          })
        } else {
          this.$Message.error('请输入必填项!')
        }
      })
    },
    openBackup (row) {
      this.backupModal = true
      this.formBackup.snapId = row.id
    },
    openMirror (row) {
      this.mirrorModal = true
      this.formMirror.snapId = row.id
    },
    crtMirror (name) {
      this.$refs[name].validate((valid) => {
        if (valid) {
          this.isdisabledFn=true
          const param = {
            id: this.formMirror.snapId,
            name: this.formMirror.name,
            desc: this.formMirror.desc
          }
          this.$Notice.info({
            title: '请求中……',
            desc: '正在生成镜像，请等待',
            duration:10
          })
          snapshotToMirror(param).then(({ data: { code, message } }) => {
            if (code !== 200) {
              this.$Notice.error({
                title: '错误',
                desc: message
              })
              this.isdisabledFn=false
              return
            }
            this.$Notice.success({
              title: '成功',
              desc: '操作成功'
            })
            this.mirrorModal=false
            this.isdisabledFn=false
            this.getData()
          }).catch(err => {
            if (err && err.response) {
              // console.error(err)
              this.$Notice.error({
                title: '错误',
                desc: err.response.data.msg
              })
              this.isdisabledFn=false
            }
          })
        } else {
          this.$Message.error('请输入必填项!')
        }
      })
    },
    crtBackup (name) {
      this.$refs[name].validate((valid) => {
        if (valid) {
          this.isdisabledFn=true
          const param = {
            id: this.formBackup.snapId,
            name: this.formBackup.name,
            desc: this.formBackup.desc
          }
          this.$Notice.info({
            title: '请求中……',
            desc: '正在生成副本，请等待',
            duration:10
          })
          snapshotToBackup(param).then(({ data: { code, message } }) => {
            if (code !== 200) {
              this.$Notice.error({
                title: '错误',
                desc: message
              })
              this.isdisabledFn=false
              return
            }
            this.$Notice.success({
              title: '成功',
              desc: '操作成功'
            })
            this.isdisabledFn=false
            this.backupModal = false
            this.getData()
          }).catch(err => {
            if (err && err.response) {
              // console.error(err)
              this.$Notice.error({
                title: '错误',
                desc: err.response.data.msg
              })
              this.isdisabledFn=false
            }
          })
        } else {
          this.$Message.error('请输入必填项!')
        }
      })
    },
    // 回滚
    rollback (row) {
      const params = Object.assign({}, {
        id: row.id,
      })
      httpRequest(this, rollback, [params]).then(res => {
        if (res.code !== 200) {
          this.$Notice.error({
            title: '错误',
            desc: res.message
          })
        } else {
          this.$Notice.success({
            title: '成功',
            desc: '请求成功'
          })
        }
        this.getData()
      })
    },

    //删除快照
    destroy (row) {
      const params = Object.assign({}, {
        id: row.id,
      })
      httpRequest(this, deleteSnapshot, [params]).then(res => {
        if (res.code !== 200) {
          this.$Notice.error({
            title: '错误',
            desc: res.message
          })
        } else {
          this.$Notice.success({
            title: '成功',
            desc: '请求成功'
          })
        }
        this.getData()
      })
    },
    cancel () {
    },
    // 搜索
    search () {
      this.currentPage = 1
      this.getData()
    },
    // 翻页
    changePage (page) {
      this.currentPage = page
      this.getData(page)
    },
    closeCreateModal () {
      this.currentTicketId = ''
      this.showCreateSheet = false
      this.modifyCreate = false
    },
    closeModal () {
      this.$emit('closeCreate')
    }
  }
}
</script>

<style lang="less">
</style>
