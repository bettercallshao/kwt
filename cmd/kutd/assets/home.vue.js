var Home = Vue.component("Home", {
  template: html`
    <div>
      <div class="list-group">
        <a
          v-for="c in channel"
          v-on:click="clickChannel(c)"
          class="list-group-item list-group-item-action"
        >
          <pre
            display="inline"
            class="mb-0"
          ><strong>Channel {{ c }}</strong></pre>
        </a>
      </div>
      <div class="list-group mt-3">
        <a v-for="m in menu" class="list-group-item list-group-item-action">
          <pre display="inline" class="mb-0">Menu {{ m }}</pre>
        </a>
      </div>
      <div class="border mt-3 p-3">
        <form v-on:submit.prevent>
          <div class="form-group row">
            <label for="name" class="col-sm-2 col-form-label">Name</label>
            <div class="col-sm-10">
              <input
                v-model="name"
                type="text"
                class="form-control"
                id="name"
                name="name"
                placeholder="kubectl"
              />
            </div>
          </div>
          <div class="form-group row">
            <label for="source" class="col-sm-2 col-form-label">Source</label>
            <div class="col-sm-10">
              <input
                v-model="source"
                type="text"
                class="form-control"
                id="source"
                name="source"
                placeholder="/home/steve/Downloads/kubectl.yaml"
              />
            </div>
          </div>
          <button
            type="submit"
            v-on:click="clickIngest()"
            class="btn btn-primary"
          >
            Ingest Menu
          </button>
        </form>
      </div>
    </div>
  `,
  data() {
    return {
      channel: [],
      menu: [],
      name: "",
      source: ""
    };
  },
  created() {
    axios.get("/channel").then(response => (this.channel = response.data.list));
    axios.get("/menu").then(response => (this.menu = response.data.list));
  },
  methods: {
    clickChannel(c) {
      this.$router.push(`channel/${c}`);
    },
    clickIngest() {
      const name = this.name;
      const source = this.source;
      axios.post("/menu", { name, source }).then(() => (window.location = "/"));
    }
  }
});
