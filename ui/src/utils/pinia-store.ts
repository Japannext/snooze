import { defineStore, storeToRefs } from "pinia";

export function defineRefStore(id, storeSetup, options?) {
  const piniaStore = defineStore(id, storeSetup, options);

  return () => {
    const usedPiniaStore = piniaStore();

    const storeRefs = storeToRefs(usedPiniaStore);

    return { ...usedPiniaStore, ...storeRefs };
  };
}
