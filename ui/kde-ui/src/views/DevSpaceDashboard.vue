<template>
    <el-table
      :data="tableData"
      style="width: 100%"
      :row-class-name="tableRowClassName"
    >
      <el-table-column prop="metadata.name" label="Name" width="180" />
      <el-table-column prop="status.link" label="Address" />
      <el-table-column prop="status.phase" label="Status" width="180" />
      <el-table-column fixed="right" label="Operations" min-width="80">
        <template #default="scope">
          <el-button
            link
            type="primary"
            size="small"
            @click.prevent="deleteRow(scope.$index)"
          >
            Remove
          </el-button>
        </template>
      </el-table-column>
    </el-table>
</template>
  
<script lang="ts" setup>
import { ref } from 'vue'

interface User {
    status: {
        phase: string
        link: string
    }
}

  const tableRowClassName = ({
    row,
    rowIndex,
  }: {
    row: User
    rowIndex: number
  }) => {
    if (rowIndex === 1) {
      return 'warning-row'
    } else if (row.status.phase === 'Ready') {
      return 'success-row'
    }
    return ''
  }
  
const tableData = ref([])
const loadData = () => {
    fetch('/api/devspace', {
        method: 'GET'
    }).then(res => {
        return res.json()
    }).then(d => {
        tableData.value = d.items
    })
}
loadData()

const deleteRow = (index: number) => {
    fetch(`/api/devspace/${tableData.value[index].metadata.name}`, {
        method: 'DELETE'
    }).finally(() => {
        loadData()
    })
}
</script>
  
<style>
.el-table .warning-row {
--el-table-tr-bg-color: var(--el-color-warning-light-9);
}
.el-table .success-row {
--el-table-tr-bg-color: var(--el-color-success-light-9);
}
</style>
