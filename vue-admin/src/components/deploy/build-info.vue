<template>
  <div>
    <Collapse v-model="showPanel">
      <Panel name="buildrecord">
        <strong style="color: #3399CC;font-size: 16px">部署客户端记录</strong>
        <div slot="content">
          <Row>
            <Table searchable border :columns="columnsTitle" :data="orders" stripe size="small"></Table>
            <div class="table-page">
              <div class="table-page-position">
                <Page :total="total" :current="currentPage" show-total show-elevator></Page>
              </div>
            </div>
          </Row>
        </div>
      </Panel>
      <Panel name="build">
        主机信息
        <p slot="content">
          <Form :label-width="90" label-position="left" class="not-validate-item">
            <Row>
              <Col span="12">
                <FormItem label="主机名">
                  <span>{{ detailData.hostname }}</span>
                </FormItem>
              </Col>
              <Col span="12">
                <FormItem label="描述">
                  <span>{{ detailData.desc }}</span>
                </FormItem>
              </Col>
              <Col span="12">
                <FormItem label="ip">
                  <span>{{ detailData.ip }}</span>
                </FormItem>
              </Col>
              <Col span="12">
                <FormItem label="端口">
                  <span>{{ detailData.port }}</span>
                </FormItem>
              </Col>
            </Row>
            <Row>
              <Col span="12">
                <FormItem label="架构">
                  <span>{{ detailData.arch }}</span>
                </FormItem>
              </Col>
              <Col span="12">
                <FormItem label="操作系统">
                  <span>{{ detailData.os }}</span>
                </FormItem>
              </Col>
              <Col span="12">
                <FormItem label="类型">
                  <span>{{ typeText }}</span>
                </FormItem>
              </Col>
              <Col span="12">
                <FormItem label="用户名">
                  <span>{{ detailData.username }}</span>
                </FormItem>
              </Col>
            </Row>
            <Row>
              <Col span="12">
                <FormItem label="创建人">
                  <span>{{ detailData.create_user }}</span>
                </FormItem>
              </Col>
              <Col span="12">
                <FormItem label="创建时间">
                  <span>{{ detailData.created_time }}</span>
                </FormItem>
              </Col>
            </Row>
          </Form>
        </p>
      </Panel>
    </Collapse>
  </div>
</template>

<script>
import { getHostDetails, getDeployRecord, getDeployStatus, choiceDeployStateColor } from '@/api/env.js'
import { httpRequest } from '@/libs/tools.js'
import { getInstanceList } from '@/api/backup'
export default {
  name: 'buildInfo',
  props: {
    bId: {
      type: [Number, String],
      required: false
    },
    typeMap: Object
  },
  data () {
    return {
      detailData: {},
      deployDetail: {},
      showDetail: false,
      showPanel: 'buildrecord',
      typeText: '',
      orders: [],
      columnsTitle: [
        {
          title: '控制台',
          key: 'serverIp',
          align: 'center'
        },
        {
          title: '安装路径',
          key: 'baseDir',
          align: 'center'
        },
        {
          title: '安装应用',
          key: 'apps',
          align: 'center',
          render: (h, params) => {
            let arr = []
            arr.push(
              h('span', { style: { cursor: 'pointer' } },
                [
                  h('span',
                    {
                      style: {
                        display: 'inline-block',
                        width: '100%',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap'
                      },
                      domProps: {
                        title: params.row.apps
                      }
                    },
                    params.row.apps)
                ]
              )

            )
            return h('span', arr)
          }
        },
        {
          title: '类型',
          width: 65,
          key: 'abTpye',
          align: 'center'
        },
        {
          title: '部署结果',
          width: 120,
          key: 'status',
          align: 'center',
          render: (h, params) => {
            let lastDeployState = getDeployStatus(params.row.status)
            let color = choiceDeployStateColor(params.row.status)
            if (!lastDeployState && !params.row.status) {
              return h('span', '')
            }

            if (lastDeployState) {
              return h('Tag', {
                props: {
                  type: 'dot',
                  color: color,
                  fade: true
                }
              }, lastDeployState)
            } else {
              return h('Tag', {
                props: {
                  type: 'dot',
                  color: color,
                  fade: true
                }
              }, params.row.status)
            }
          }
        },
        {
          title: '创建时间',
          width: 100,
          key: 'created_time',
          align: 'center'
        },
        {
          title: '更新时间',
          width: 100,
          key: 'updated_time',
          align: 'center'
        }
      ],
      total: 0,
      currentPage: 1
    }
  },
  computed: {
    deployColor () {
      return '#515a6e'
      // if (!this.detailData.deploy) return '#515a6e'
      // let color = choiceDeployStateColor(this.detailData.deploy.status_text)
      // return color
    }
  },
  methods: {
    detail (_id) {
      this.showDetail = true
      httpRequest(this, getHostDetails, [_id]).then(res => {
        this.detailData = res.response
        this.typeText = this.typeMap[res.response.type]
        // this.deployDetail = res.deploy
      })
    },
    getRecordData (page = 1, _id) {
      let hostId = _id ? _id : this.bId
      const params = Object.assign({}, {
        id: hostId
      })

      httpRequest(this, getDeployRecord, [page, params]).then(res => {
        if (res.code !== 200) {
          this.$Notice.error({
            title: '错误',
            desc: res.message
          })
          return
        }
        if (res.response.data) {
          this.orders = res.response.data
          this.total = res.response.totalNum
        } else {
          this.orders = []
          this.total = 0
        }
        this.currentPage = page
      })
    }
  }
}
</script>
