<template>
<MenuComponent :userinfo=userinfo :viewname=viewname :token=token />
Scope:
<select @change='get_assets()' v-model='chosen_scope'>
  <option v-for='(scope, id) in scopes' :value='scope.ID' :key=id>{{ scope.Name }}</option>
</select>
<table id='datatable'>
<thead>
  <tr>
    <th>Big Category</th>
    <th>Small Category</th>
    <th>Name</th>
    <th>Owner</th>
    <th>C</th>
    <th>I</th>
    <th>A</th>
    <th>Threat</th>
    <th>Vulnerability</th>
    <th>Control</th>
    <th>Possibility</th>
    <th>Impact</th>
    <th>Risk</th>
    <th>Actions</th>
  </tr>
</thead>
<tbody>
  <template v-for='asset in assets'>
    <tr v-for='(risk, risk_i) in asset.Risks' :key=risk_i>
      <td class="asset" v-if='!risk_i' :rowspan='asset.Risks.length'>{{ asset.BigCategory }}</td>
      <td class="asset" v-if='!risk_i' :rowspan='asset.Risks.length'>{{ asset.SmallCategory }}</td>
      <td class="asset" v-if='!risk_i' :rowspan='asset.Risks.length'>{{ asset.Name }}</td>
      <td class="asset" v-if='!risk_i' :rowspan='asset.Risks.length'>{{ asset.Owner }}</td>
      <td class="asset" v-if='!risk_i' :rowspan='asset.Risks.length'>{{ asset.Value.Confidentiality }}</td>
      <td class="asset" v-if='!risk_i' :rowspan='asset.Risks.length'>{{ asset.Value.Integrity }}</td>
      <td class="asset" v-if='!risk_i' :rowspan='asset.Risks.length'>{{ asset.Value.Availability }}</td>
      <template v-if='risk.mode == "read"'>
      <td>{{ risk.Threat }}</td>
      <td>{{ risk.Vulnerability }}</td>
      <td>{{ risk.CurrentControl }}</td>
      <td>{{ risk.Possibility }}</td>
      <td>{{ risk.Impact }}</td>
      <td>{{ calculate_risk(
               asset.Value.Confidentiality,
               asset.Value.Integrity,
               asset.Value.Availability,
               risk.Possibility,
               risk.Impact
           ) }}</td>
      <td>
        <button @click='edit_risk(risk)'>Edit</button>
        <button v-if='risk_i + 1 == asset.Risks.length' @click='add_risk(asset)'>Add</button>
        <button @click='delete_risk(asset, risk)'>Delete</button>
      </td>
      </template>
      <template v-else>
      <td><input v-model.lazy='risk.Threat' /></td>
      <td><input v-model.lazy='risk.Vulnerability' /></td>
      <td><input v-model.lazy='risk.CurrentControl' /></td>
      <td>
        <select v-model.number='risk.Possibility'>
          <option value='' selected disabled>Choose</option>
          <option v-for='(option, id) in config.Possibility' :value='option.value' :key=id>{{ option.text }}</option>
        </select>
      </td>
      <td>
        <select v-model.number='risk.Impact'>
          <option value='' selected disabled>Choose</option>
          <option v-for='(option, id) in config.Impact' :value='option.value' :key=id>{{ option.text }}</option>
        </select>
      </td>
      <td>{{ calculate_risk(
               asset.Value.Confidentiality,
               asset.Value.Integrity,
               asset.Value.Availability,
               risk.Possibility,
               risk.Impact
           ) }}</td>
      <td>
         <button @click='update_asset(asset)'>{{ risk.mode == "edit" ? "Update" : "Add" }}</button>
      </td>
      </template>
    </tr>
  </template>
</tbody>
</table>
</template>

<script>
import json from '@/assets/config.json'
import MenuComponent from '@/components/Menu.vue'

export default {
  components: {
    MenuComponent
  },
  data() {
    return {
      viewname: 'RiskAssessmentView',
      userinfo: {},
      scopes: [],
      chosen_scope: '',
      assets: [],
      config: json,
    }
  },
  beforeMount() {
    this.get_scopes();
  },
  methods: {
    get_scopes: function () {
      let ref = this;
      fetch('/api/getscopesbyuser', {
        method: 'get',
      }).then((response) => {
        if (!response.ok) throw new Error(response.statusText)
        ref.token = response.headers.get('X-Csrf-Token');
        return response.json();
      }).then((res) => {
        ref.userinfo = res.UserInfo;
        ref.scopes = res.Scopes;
        ref.chosen_scope = res.Scopes[0].ID;
        ref.get_assets();
      }).catch(function(err) {
        console.log(err);
      });
    },
    to_int: function (val) {
      if (typeof val === 'number')
        return val;
      if (parseInt(val) === 'number')
        return parseInt(val);
      return 0;
    },
    calculate_risk: function (c, i, a, possibility, impact) {
      let c_int = this.to_int(c);
      let i_int = this.to_int(i);
      let a_int = this.to_int(a);
      let possibility_int = this.to_int(possibility);
      let impact_int = this.to_int(impact);
      return (c_int + i_int + a_int) * possibility_int * impact_int;
    },
    get_assets: function () {
      let ref = this;
      if (!ref.chosen_scope.length) {
        alert("Please choose a scope")
        return;
      }
      fetch('/api/getassets/' + ref.chosen_scope, {
        method: 'get',
      }).then((response) => {
        if (!response.ok) throw new Error(response.statusText)
        return response.json();
      }).then((res) => {
        ref.userinfo = res.UserInfo;

        if (res.Assets == null) {
          ref.assets = [];
        }
        else {
          ref.assets = res.Assets;
          ref.assets.forEach((ast) => {
            if (ast.Risks.length == 0) {
              ast.Risks.push({
                Threat: "",
                Vulnerability: "",
                CurrentControl: "",
                Possibility: "",
                Integrity: "",
                mode: "add",
              });
            } else {
              ast.Risks.forEach((risk) => risk.mode = "read");
            }
          });
        }
      }).catch(function(err) {
        console.log(err);
      });
    },
    add_risk: function (asset) {
      asset.Risks.push({
        Threat: "",
        Vulnerability: "",
        CurrentControl: "",
        Possibility: "",
        Integrity: "",
        mode: "add",
      });
    },
    edit_risk: function (risk) {
      risk.mode = 'edit'
    },
    update_asset (asset) {
      let ref = this;
      fetch('/api/updateasset', {
        method: 'post',
        headers: new Headers({
          'Content-Type': 'application/json',
          'X-Csrf-Token': ref.token,
        }),
        body: JSON.stringify(asset)
      }).then(function(response) {
        if (!response.ok) throw new Error(response.statusText)
        ref.get_assets();
      }).catch(function(err) {
        console.log(err);
      });
    },
    delete_risk: function (asset, risk) {
      asset.Risks.splice(asset.Risks.indexOf(risk), 1);
      this.update_asset(asset);
    }
  }
};
</script>
