var Channel = Vue.component("Channel", {
  template: html`
    <div>
      <p>
        <strong>{{ errors.join(" ") || menu.name }}</strong>
      </p>
      <p>Channel {{ $route.params.channel}}</p>
      <p>{{ menu.help }}</p>
      <p>{{ menu.version }} {{ menu.hash }}</p>
      <router-view class="mt-5" :menu="menu" :token="token" :socket="socket"></router-view>
      <div class="custom-control custom-switch mt-3">
        <input type="checkbox" class="custom-control-input" id="markdown" v-model="checked">
        <label class="custom-control-label" for="markdown">Markdown</label>
      </div>
      <div class="border p-2">
        <div v-if="checked" v-html="latestMarked"></div>
        <div v-else><pre>{{ latestRaw }}</pre></div>
      </div>
      <div class="custom-control custom-switch mt-3">
        <input type="checkbox" class="custom-control-input" id="showLog" v-model="showLog">
        <label class="custom-control-label" for="showLog">Show log</label>
      </div>
      <div v-if="showLog" class="border p-2 overflow-auto" style="height: 360px">
        <pre>{{ logRaw }}</pre>
      </div>
    </div>
  `,
  data() {
    return {
      menu: {},
      token: "",
      errors: [],
      socket: {},
      outs: [],
      latest: [],
      checked: false,
      showLog: false,
    };
  },
  computed: {
    latestRaw() {
      return this.latest.join("\n")
    },
    latestMarked() {
      return marked(this.latestRaw, { sanitize: true });
    },
    logRaw() {
      return this.outs.join("\n")
    },
  },
  created() {
    const channel = this.$route.params.channel;
    axios
      .get(`/channel/${channel}`)
      .then(response => {
        this.menu = response.data.menu;
        this.token = response.data.token;
      })
      .catch(error => {
        this.errors.push(error.response.data.error);
      });
    let prefix = "ws";
    if (document.location.protocol === "https:") {
      prefix = "wss";
    }
    let host = document.location.host;
    let wsUrl = `${prefix}://${host}/ws/front/${channel}`;
    this.socket = new WebSocket(wsUrl);
    this.socket.onmessage = event => {
      try {
        data = JSON.parse(event.data);
        if (data.token == this.token) {
          this.outs.push(data.payload.text);
          if (data.payload.pipe == 'stdout') {
            this.latest.push(data.payload.text);
          } else if (data.payload.pipe == 'echo') {
            this.latest = [];
          }
        }
      } catch (_) {}
    };
    this.socket.onclose = event => {
      this.errors.push("Connection refused.");
    };
  }
});
