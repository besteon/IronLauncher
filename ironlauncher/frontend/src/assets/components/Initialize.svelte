<script>
	import ironlogo from '../images/iron.png'
    import launcherlogo from '../images/launcher.png'
    import {StartUp} from '../../../wailsjs/go/main/App.js'

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
    <img alt="Iron logo" id="logo" src="{ironlogo}">
    <img alt="Launcher logo" id="logo" src="{launcherlogo}">

    
</main>

<style>

	#logo {
        display: block;
        width: 50%;
        height: 50%;
        margin: auto;
        padding: 10% 0 0;
        background-position: center;
        background-repeat: no-repeat;
        background-size: 100% 100%;
        background-origin: content-box;
    }

</style>