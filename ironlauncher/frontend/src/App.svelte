<script>
  import Installer from './views/Installer.svelte';
  import Initialize from './views/Initialize.svelte';
  import Launcher from './views/Launcher.svelte';
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
    AreDepsInstalled().then(result => {
      if (result == true) {
        window.location.hash = '#initialize'
      } else {
        window.location.hash = '#installer'
      }
    })
  }

  areDepsInstalled();

</script>

<svelte:window on:hashchange={routeChange} />

<!-- <nav>
  <a href="/#installer" class:active={page === Installer}>Installer</a>
  <a href="/#initialize" class:active={page === Initialize}>Initialize</a>
  <a href="/#launcher" class:active={page === Launcher}>Launcher</a>
</nav> -->

<main>
  <svelte:component this={page} />
</main>

<style>
  :global(body) {
    background-color: rgb(138, 181, 216);
  }
</style>
