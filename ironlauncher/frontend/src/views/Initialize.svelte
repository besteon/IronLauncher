<script>
	import ironlogo from '../assets/iron.png'
    import launcherlogo from '../assets/launcher.png'
    import {StartUp} from '../../wailsjs/go/main/App.js'

    (function wait_for_install() {
        StartUp().then(result => {
            console.log(result)
            if (result) {
                window.location.hash = '#launcher'
            } else {
                setTimeout(wait_for_install, 1000)
            }
        })
    })()
</script>

<main>
    <div class="window">
        <div class="logoContainer">
            <img alt="Iron logo" class="logo" src="{ironlogo}">
            <img alt="Launcher logo" class="logo" src="{launcherlogo}">
            <div class="padding"></div>
        </div>
    </div>

</main>
  
<style>

    .window {
        background-color: rgb(138, 181, 216);
        width: 100%;
        height: 100%;
    }

    .padding {
        margin-bottom: 2%;
    }

    .window > * {
        text-align: center;
        margin: auto;
        min-height: 100%;
        min-width: 30%;
    }

    .logo {
        width: 50%;
        margin-left: auto;
        margin-right: auto;
    }

</style>