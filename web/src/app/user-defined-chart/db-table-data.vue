<template>
  <div>
    <h1>选择数据库和表</h1>
    <el-row>
      <el-form ref="getTableForm" :inline="true" :model="dbform" label-width="120px">
        <el-form-item label="数据库名">
          <el-select @change="getTable" v-model="dbform.dbName" filterable placeholder="请选择数据库">
            <el-option
              v-for="item in dbOptions"
              :key="item.database"
              :label="item.database"
              :value="item.database"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="表名">
          <el-select
            @change="getColumn"
            v-model="dbform.tableName"
            :disabled="!dbform.dbName"
            filterable
            placeholder="请选择表"
          >
            <el-option
              v-for="item in tableOptions"
              :key="item.tableName"
              :label="item.tableName"
              :value="item.tableName"
            ></el-option>
          </el-select>
        </el-form-item>
      </el-form>
    </el-row>

    <h1>查询字段，第一个字段映射为X轴/半径轴，第二个字段映射为Y轴/角度轴</h1>
    <el-row>
      <el-col :span="24">
        <el-select style="width:100%" placeholder="查询字段" v-model="queryColumns" multiple clearable>
          <el-option
            v-for="item in tableColumnsOption"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          ></el-option>
        </el-select>
      </el-col>
    </el-row>

    <!-- <h1>选择数据维度</h1>
    <el-row>
      <el-form :inline="true">
        <el-form-item label="x">
          <el-select placeholder="X轴/半径轴" v-model="dataDimension.x">
            <el-option
              v-for="item in tableColumnsOption"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="y">
          <el-select placeholder="Y轴/角度轴" v-model="dataDimension.y">
            <el-option
              v-for="item in tableColumnsOption"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="z">
          <el-select placeholder v-model="dataDimension.z">
            <el-option
              v-for="item in tableColumnsOption"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
      </el-form>
    </el-row>-->

    <h1>筛选条件</h1>
    <el-row>
      <el-button type="primary" @click="querySettings.push({ uuid: uuid(), data: {} });">添加</el-button>
      <el-button type="danger" @click="querySettings=[]">清空</el-button>
      <!-- <el-button type="success" @click="applyQuerySetting">应用</el-button> -->
    </el-row>
    <el-row>
      <el-row v-for="(item, index) in querySettings" :key="item.uuid">
        <el-col :span="22">
          <query-setting
            :fields="tableColumnsOption"
            v-on:update="(val) => {querySettings[index].data=val}"
          ></query-setting>
        </el-col>
        <el-col :span="2">
          <el-button @click="querySettings.splice(index, 1)">删除</el-button>
        </el-col>
      </el-row>
    </el-row>

    <el-row>
      <el-col :span="2" :offset="22">
        <el-button @click="getData();" type="primary">确 定</el-button>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import service from "@/utils/request";
import uuid from "uuid/v1";
import QuerySetting from "@/components/query-setting";
import { getDB, getTable, getColumn } from "@/api/autoCode.js";

export default {
  name: "DbTableData",
  components: { QuerySetting },
  data() {
    return {
      //////////////////
      tableColumnsOption: [],
      querySettings: [],
      //////////////////
      dbform: {
        dbName: "",
        tableName: ""
      },
      dbOptions: [],
      tableOptions: [],
      tableColumns: [],
      //////////////////
      queryColumns: []
      //////////////////
    };
  },
  methods: {
    ///////////////////////////////////
    uuid() {
      uuid();
    },
    ///////////////////////////////////
    async getDb() {
      const res = await getDB();
      if (res.code == 0) {
        this.dbOptions = res.data.dbs;
      }
    },
    async getTable() {
      const res = await getTable({ dbName: this.dbform.dbName });
      if (res.code == 0) {
        this.tableOptions = res.data.tables;
      }
      this.dbform.tableName = "";
    },
    async getColumn() {
      const res = await getColumn(this.dbform);
      if (res.code == 0) {
        if (res.data.columns) {
          this.tableColumns = res.data.columns.map(item => {
            return item.columnName;
          });

          this.tableColumnsOption = res.data.columns.map(item => {
            return {
              label:
                item.columnComment != "" ? item.columnComment : item.columnName,
              value: item.columnName
            };
          });
        }
      }
    },
    ///////////////////////////////////
    async getData() {
      // console.log(this.querySettings)
      var where = "";
      for (let i = 0; i < this.querySettings.length; i++) {
        const e = this.querySettings[i];
        if (!e.data || !e.data.field) {
          continue;
        }

        var s = "`[field]` [mode] '[value]'";
        s = s.replace("[field]", e.data.field);
        s = s.replace("[mode]", e.data.mode);
        let v = "";
        if (e.data.mode === "like") {
          v = "%" + e.data.value + "%";
        } else {
          v = e.data.value;
        }
        s = s.replace("[value]", v);

        if (i != this.querySettings.length - 1) {
          where += s + " and ";
        } else {
          where += s;
        }
      }

      // 查询参数
      var params = {
        dbName: this.dbform.dbName,
        tableName: this.dbform.tableName,
        queryColumns: this.queryColumns,
        filter: where
      };
      // console.log(params);

      const res = await service({
        url: "/dbTableData/data",
        method: "get",
        params
      });
      if (res.code == 0) {
        // console.log(res.data)
        if (res.data) {
          this.$emit("data", {
            dbName: this.dbform.dbName,
            tableName: this.dbform.tableName,
            queryColumns: this.queryColumns,
            filter: where,
            data: res.data
          });
        }

        this.$message({
          type: "success",
          message: "getData成功"
        });
      }
    }
  },
  created() {
    this.getDb();
  }
};
</script>

<style lang="scss" scoped>
</style>