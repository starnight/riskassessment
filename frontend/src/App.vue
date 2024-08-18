<template>
  <h1>Asset Risk Management</h1>
  <router-view/>
</template>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#datatable {
  font-family: Arial, Helvetica, sans-serif;
  border-collapse: collapse;
  width: 100%;
}

#datatable td, #datatable th {
  border: 1px solid #ddd;
  padding: 8px;
}

#datatable tr:nth-child(even){background-color: #f2f2f2;}

#datatable td.asset {background-color: #f2f2f2;}

#datatable tr:hover {background-color: #ddd;}

#datatable th {
  padding-top: 12px;
  padding-bottom: 12px;
  text-align: left;
  background-color: #04AA6D;
  color: white;
}
</style>
