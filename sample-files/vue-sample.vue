<template>
  <div>
    <q-form @submit="onSubmit">
      <q-btn type="submit">Submit</q-btn>
    </q-form>
    <q-dialog v-model="showDialog">
      <div>Dialog content</div>
    </q-dialog>
  </div>
</template>

<script>
export default {
  name: 'VueSample',
  data() {
    return {
      showDialog: false
    }
  },
  methods: {
    onSubmit() {
      console.log('Form submitted')
    }
  }
}
</script>
