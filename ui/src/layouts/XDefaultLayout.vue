<script setup lang="ts"></script>

<template>
  <div class="absolute"></div>
  <XMenu />
  <div class="h-100">
    <div class="layout-main-container">
      <router-view />
    </div>
  </div>
</template>
>

<style>
body {
  font-family: "Lato", sans-serif;
  margin: 0;
  padding: 0;
  min-height: 100%;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.layout-wrapper {
  min-height: 100vh;
}

.layout-top-bar {
  position: absolute;
  height: 58px;
  padding: 3px;
}

.layout-main-container {
  top: 58px;
  bottom: 58px;
  padding: 10px;
  position: absolute;
}

.layout-main {
  flex: 1 1 auto;
  padding-bottom: 2rem;
}

.layout-topbar {
  .layout-topbar-logo-container {
    width: 20rem;
    display: flex;
    align-items: center;
  }

  .layout-topbar-logo {
    display: inline-flex;
    align-items: center;
    font-size: 1.5rem;
    border-radius: var(--content-border-radius);
    font-weight: 500;
    gap: 0.5rem;

    svg {
      width: 3rem;
    }
  }

  .layout-topbar-action {
    display: inline-flex;
    justify-content: center;
    align-items: center;
    border-radius: 50%;
    width: 2.5rem;
    height: 2.5rem;
    cursor: pointer;

    i {
      font-size: 1.25rem;
    }

    span {
      font-size: 1rem;
      display: none;
    }
  }

  .layout-menu-button {
    margin-right: 0.5rem;
  }

  .layout-topbar-menu-button {
    display: none;
  }

  .layout-topbar-actions {
    margin-left: auto;
    display: flex;
    gap: 1rem;
  }

  .layout-topbar-menu-content {
    display: flex;
    gap: 1rem;
  }

  .layout-config-menu {
    display: flex;
    gap: 1rem;
  }
}
</style>
