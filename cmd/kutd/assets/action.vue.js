var Action = Vue.component("Action", {
  template: html`
    <div>
      <p><strong>{{ action.name }}</strong></p>
      <p>{{ action.help }}</p>
      <form v-on:submit.prevent>
        <div v-for="p in action.params" class="form-group row">
          <label :for="p.name" class="col-sm-2 col-form-label"
            >{{ p.name }}</label
          >
          <div class="col-sm-10">
            <input
              v-model="form[p.name]"
              :id="p.name"
              :name="p.name"
              :placeholder="p.help"
              type="text"
              class="form-control-plaintext border p-2"
            />
          </div>
        </div>
        <button
          type="submit"
          v-on:click="clickExecute()"
          class="btn btn-primary"
        >
          Execute
        </button>
      </form>
    </div>
  `,
  props: ["menu", "token", "socket"],
  computed: {
    action() {
      if (this.menu.actions) {
        const name = this.$route.params.action;
        return this.menu.actions.filter(a => {
          return a.name == name;
        })[0];
      } else {
        return {};
      }
    },
    form() {
      return this.action.params.reduce((form, p) => {
        form[p.name] = p.value;
        return form;
      }, {});
    }
  },
  methods: {
    clickExecute() {
      command = {
        token: this.token,
        action: {
          name: this.action.name,
          template: this.action.template,
          params: Object.keys(this.form).map(key => ({
            name: key,
            value: this.form[key]
          }))
        }
      };
      this.socket.send(JSON.stringify(command));
    }
  }
});
