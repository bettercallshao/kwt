var Menu = Vue.component("Menu", {
  template: html`
    <div class="list-group">
      <a
        v-for="a in menu.actions"
        v-on:click="clickAction(a.name)"
        class="list-group-item list-group-item-action"
      >
        <pre
          display="inline"
          class="mb-0"
        ><strong>{{ a.name.padEnd(30) }}</strong>{{ a.help }}</pre>
      </a>
    </div>
  `,
  props: ["menu"],
  methods: {
    clickAction(name) {
      this.$router.push(this.$route.path.trimRight("/") + `/action/${name}`);
    }
  }
});
