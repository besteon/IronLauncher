<script>
  import Installer from './assets/components/Installer.svelte';
  import Initialize from './assets/components/Initialize.svelte';
  import Launcher from './assets/components/Launcher.svelte';
  import {AreDepsInstalled} from '../wailsjs/go/main/App.js'

  const routingMap ={
    '#installer': Installer,
    '#initialize': Initialize,
    '#launcher': Launcher
  }

  let page;

  function routeChange() {
    page = routingMap[location.hash] || Launcher
  }

  function areDepsInstalled(){
    console.log('test')
    AreDepsInstalled().then(result => {
      console.log(result)
      if (result == true) {
        window.location.hash = '#launcher'
      } else {
        window.location.hash = '#installer'
      }
    })
  }

  areDepsInstalled();

</script>

<svelte:window on:hashchange={routeChange} />

<nav>
  <a href="/#installer" class:active={page === Installer}>Installer</a>
  <a href="/#initialize" class:active={page === Initialize}>Initialize</a>
  <a href="/#launcher" class:active={page === Launcher}>Launcher</a>
</nav>

<main>
  <svelte:component this={page} />
</main>

<style>

</style>
