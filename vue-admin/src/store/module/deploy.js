export default {
  state: {
    // SOR发布也多个页面同步使用TicketID
    ticketId: ''
  },
  mutations: {
    setTicketId (state, params) {
      state.ticketId = params
    }
  }
}
