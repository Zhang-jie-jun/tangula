export default {
  bind (el, binding, vnode) {
    const { value } = binding
    const { permits, control } = value
    const pass = control.every((d) => {
      return permits.includes(d)
    })
    if (!permits.length) {
      return
    }
    if (!pass) {
      el.style.display = 'none'
    }
  },
  componentUpdated (el, binding, vnode) {
    const { value } = binding
    const { permits, control } = value
    const pass = control.every((d) => {
      return permits.includes(d)
    })
    // console.log(pass)
    if (!permits.length) {
      return
    }
    if (!pass) {
      el.style.display = 'none'
    }
  }
}
