<script setup></script>

<template>
  <h1>仪表板</h1>
</template>

<style lang="scss" scoped>
.page {
  @apply p-0;
  .my-card-box {
    @apply p-4;
    & + .my-card-box {
      @apply pt-0;
    }
  }

  .my-card {
    @apply box-border bg-white rounded h-auto px-6 py-8 overflow-hidden shadow-sm;
    .my-card-title {
      @apply pb-5 border-t-0 border-l-0 border-r-0 border-b border-solid border-gray-100;
    }
  }

  .my-top-card {
    @apply h-72 flex items-center justify-between text-gray-500;
    &-left {
      @apply h-full flex flex-col w-auto;
      &-title {
        @apply text-3xl text-gray-600;
      }

      &-dot {
        @apply mt-4 text-gray-600 text-lg;
      }

      &-item {
        + .my-top-card-left-item {
          margin-top: 24px;
        }

        margin-top: 14px;
      }
    }

    &-right {
      height: 600px;
      width: 600px;
      margin-top: 28px;
    }
  }

  :deep(.el-card__header) {
    @apply p-0  border-gray-200;
  }

  .card-header {
    @apply pb-5 border-b border-solid border-gray-200 border-t-0 border-l-0 border-r-0;
  }

  .quick-entrance-items {
    @apply flex items-center justify-center text-center text-gray-800;
    .quick-entrance-item {
      @apply px-8 py-6 flex items-center flex-col transition-all duration-100 ease-in-out rounded-lg cursor-pointer;
      &:hover {
        @apply shadow-lg;
      }

      &-icon {
        @apply flex items-center h-16 w-16 rounded-lg justify-center mx-0 my-auto text-2xl;
      }

      p {
        @apply mt-2.5;
      }
    }
  }
}

.dashboard-icon {
  @apply flex items-center text-xl mr-2 text-blue-400;
}
</style>
