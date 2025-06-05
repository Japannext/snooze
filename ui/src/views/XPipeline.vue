<script setup lang="ts"></script>

<template>
  <div>
    <Tabs value="groupings">
      <TabList>
        <Tab value="groupings">Groupings</Tab>
        <Tab value="mappings">Mappings</Tab>
      </TabList>
      <TabPanels>
        <TabPanel value="groupings"><XGroupingTab /></TabPanel>
        <TabPanel value="mappings"><XMappingTab /></TabPanel>
      </TabPanels>
    </Tabs>
  </div>
</template>
