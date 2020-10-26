var Channel = Vue.component("Channel", {
  template: html`
    <div>
      <p>
        <strong
          >{{ errors.join(" ") || menu.name }} ({{ $route.params.channel
          }})</strong
        >
      </p>
      <router-view :menu="menu" :token="token" :socket="socket"></router-view>
      <p class="mt-5">Latest</p>
      <div class="border p-2" v-html="compiled"></div>
      <p class="mt-5">Output</p>
      <div class="border p-2">
        <pre>{{ outs.join("\\n") }}</pre>
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
      latest: []
    };
  },
  computed: {
    compiled: function() {
      return marked(this.latest.join("\n"), { sanitize: true });
    }
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
