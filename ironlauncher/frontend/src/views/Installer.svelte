<script>
	import ironlogo from '../assets/iron.png'
    import launcherlogo from '../assets/launcher.png'
    import {AreDepsInstalled, InstallDependencies} from '../../wailsjs/go/main/App.js'

	function installDeps() {
		InstallDependencies().then(result => {
            (function wait_for_install() {
                AreDepsInstalled().then(result => {
                    if (result) {
                        window.location.hash = '#launcher'
                    } else {
                        setTimeout(wait_for_install, 1000)
                    }
                })
            })()
        })
	}
</script>

<main>
    <div class="window">
        <div class="logoContainer">
            <img alt="Iron logo" class="logo" src="{ironlogo}">
            <img alt="Launcher logo" class="logo" src="{launcherlogo}">
            <div class="padding"></div>
        </div>
        <div class="install">
            <button class="btn" on:click={installDeps}>Install</button>
        </div>
    </div>

</main>
  
<style>

    .window {
        background-color: rgb(138, 181, 216);
        width: 100%;
        height: 100%;
    }

    .install {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
    }

    .install > * {
        margin: auto;
        width: 100%;
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

    .btn {
        background-image: linear-gradient(#f7f8fa ,#e7e9ec);
        border-color: #adb1b8 #a2a6ac #8d9096;
        border-style: solid;
        border-width: 1px;
        border-radius: 3px;
        box-shadow: rgba(255,255,255,.6) 0 1px 0 inset;
        box-sizing: border-box;
        color: #0f1111;
        cursor: pointer;
        display: inline-block;
        font-family: "Amazon Ember",Arial,sans-serif;
        font-size: 14px;
        height: 29px;
        font-size: 13px;
        outline: 0;
        overflow: hidden;
        padding: 0 11px;
        text-align: center;
        text-decoration: none;
        text-overflow: ellipsis;
        user-select: none;
        -webkit-user-select: none;
        touch-action: manipulation;
        white-space: nowrap;
    }
    .btn:active {
        border-bottom-color: #a2a6ac;
    }
    .btn:active:hover {
        border-bottom-color: #a2a6ac;
    }
    .btn:hover {
        border-color: #a2a6ac #979aa1 #82858a;
    }
    .btn:focus {
        border-color: #e77600;
        box-shadow: rgba(228, 121, 17, .5) 0 0 3px 2px;
        outline: 0;
    }

    .logo {
        width: 50%;
        margin-left: auto;
        margin-right: auto;
    }

</style>