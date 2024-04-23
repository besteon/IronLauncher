<script>
	import ironlogo from '../images/iron.png'
    import launcherlogo from '../images/launcher.png'
    import {InstallDependencies} from '../../../wailsjs/go/main/App.js'

	function installDeps() {
		InstallDependencies().then(result => {
            if (result) {
                window.location.hash = '#launcher'
            }
        })
	}
</script>

<main>
    <img alt="Iron logo" id="logo" src="{ironlogo}">
    <img alt="Launcher logo" id="logo" src="{launcherlogo}">

    <div class="roms-box" id="roms-button">
        <button class="btn" on:click={installDeps}>Install</button>
    </div>
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

	.roms-box .btn {
        width: 60px;
        height: 30px;
        line-height: 30px;
        border-radius: 3px;
        border: none;
        margin: 0 0 0 20px;
        padding: 0 8px;
        cursor: pointer;
    }
</style>