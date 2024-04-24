<script>
    import ironlogo from '../images/iron.png'
    import launcherlogo from '../images/launcher.png'
    import {WindowHide, WindowShow, EventsOn} from '../../../wailsjs/runtime/runtime.js'
    import {GetRomsFolder, GetRoms, Play} from '../../../wailsjs/go/main/App.js'

	let romsFolder = "";
    let roms = [];

    function getRomsFolder() {
        GetRomsFolder().then(result => {
            romsFolder = result
            getRoms(romsFolder)
        })
    }

    function getRoms(folder) {
        roms = []
        GetRoms(folder).then(result => {
            console.log(result)
            for (let i = 0; i < result.length; i++) {
                console.log(result[i])
                roms = roms.concat({
                    value: result[i]
                })
            } 
        })
    }

    function play() {
        Play(romsFolder, "TODO").then(result => {
            WindowHide()
        })
    }

    EventsOn("EMULATOR_CLOSED", function() {
        WindowShow()
    })
</script>

<main>
    <img alt="Iron logo" id="logo" src="{ironlogo}">
    <img alt="Launcher logo" id="logo" src="{launcherlogo}">

    <div class="roms-box" id="roms-button">
        <input autocomplete="off" bind:value={romsFolder} class="roms" id="name" type="text"/>
        <button class="btn" on:click={getRomsFolder}>ROMs</button>
    </div>

    {#each roms as rom}
        <div>
            <input class="roms" type="text" bind:value={rom.value}/>
        </div>
    {/each}

    <button class="btn" on:click={play}>Play!</button>
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

    .roms-box .btn:hover {
        background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
        color: #333333;
    }

    .roms-box .roms {
        border: none;
        border-radius: 3px;
        outline: none;
        height: 30px;
        width: 400px;
        line-height: 30px;
        padding: 0 10px;
        background-color: rgba(240, 240, 240, 1);
        -webkit-font-smoothing: antialiased;
    }

    .roms-box .roms:hover {
        border: none;
        background-color: rgba(255, 255, 255, 1);
    }

    .roms-box .roms:focus {
        border: none;
        background-color: rgba(255, 255, 255, 1);
    }

</style>