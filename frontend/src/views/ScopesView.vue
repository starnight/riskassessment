<template>
<MenuComponent :userinfo=userinfo :viewname=viewname :token=token />
<table id="datatable">
<thead>
  <tr>
    <th>Scope Name</th>
    <th>Actions</th>
  </tr>
</thead>
<tbody>
  <tr v-for='(scope, id) in scopes' :key=id>
  <template v-if='scope.mode == "read"'>
    <td>{{ scope.Name }}</td>
    <td>
      <button @click='edit_scope(scope)'>Edit</button>
    </td>
  </template>
  <template v-else>
    <td><input v-model.lazy='scope.Name' /></td>
    <td>
      <button @click='update_scope(scope)' v-if='scope.mode == "edit"'>Update</button>
      <button @click='add_scope(scope)' v-else>Add</button>
    </td>
  </template>
  </tr>
</tbody>
</table>
</template>

<script>
import MenuComponent from '@/components/Menu.vue'

export default {
  components: {
    MenuComponent
  },
  data() {
    return {
      viewname: 'ScopesView',
      userinfo: {},
      scopes: [],
      token: '',
    }
  },
  beforeMount() {
    this.get_scopes();
  },
  methods: {
    get_scopes: function () {
      let ref = this;
      fetch('/api/getscopes', {
        method: 'get',
      }).then((response) => {
        if (!response.ok) throw new Error(response.statusText)
        ref.token = response.headers.get('X-Csrf-Token');
        return response.json();
      }).then((res) => {
        if (res.UserInfo.Role != 1) throw new Error("Only Administrator available")
        ref.userinfo = res.UserInfo;

        if (res.Scopes == null) {
          ref.scopes = [];
        }
        else {
          ref.scopes = res.Scopes;
          ref.scopes.forEach((ast) => {
            ast.mode = 'read';
          });
        }
        ref.scopes.push({
          Name: '',
          mode: 'add',
        });
      }).catch(function(err) {
        console.log(err);
        window.location.replace('/assets/assets.html');
      });
    },
    add_scope: function (scope) {
      let ref = this;
      fetch('/api/addscope', {
        method: 'post',
        headers: new Headers({
          'Content-Type': 'application/json',
          'X-Csrf-Token': ref.token,
        }),
        body: JSON.stringify(scope)
      }).then(function(response) {
        if (!response.ok) throw new Error(response.statusText)
        ref.get_scopes();
      }).catch(function(err) {
        console.log(err);
      });
    },
    edit_scope: function (scope) {
      scope.mode = 'edit'
    },
    update_scope (scope) {
      let ref = this;
      fetch('/api/updatescope', {
        method: 'post',
        headers: new Headers({
          'Content-Type': 'application/json',
          'X-Csrf-Token': ref.token,
        }),
        body: JSON.stringify(scope)
      }).then(function(response) {
        if (!response.ok) throw new Error(response.statusText)
        ref.get_scopes();
      }).catch(function(err) {
        console.log(err);
      });
    },
  }
};
</script>
