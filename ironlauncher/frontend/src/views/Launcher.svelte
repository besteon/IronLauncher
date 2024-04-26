<script>
    import ironlogo from '../assets/iron.png'
    import launcherlogo from '../assets/launcher.png'
    import {WindowHide, WindowShow, EventsOn} from '../../wailsjs/runtime/runtime.js'
    import {GetRomsFolder, GetRoms, Play, SaveDefaults, GetSettings} from '../../wailsjs/go/main/App.js'
    import { Label } from "$lib/components/ui/label";
    import { Checkbox } from "$lib/components/ui/checkbox";

	let romsFolder = "";
    let selectedRom = "-- GAME --";
    let selectedMode = "Vanilla";
    let roms = [];

    let modes = [
        { value: "Vanilla" },
        { value: "Classic Randomizer" },
        { value: "Classic IronMon" },
        { value: "Doubles IronMon" },
        { value: "Ultimate IronMon" },
        { value: "Survival IronMon" },
        { value: "Kaizo IronMon" },
        { value: "Super Kaizo IronMon" },
    ];

    let settings = [
        { value: "Setting 1" },
        { value: "Setting 2" },
        { value: "Setting 3" },
    ];

    function getRomsFolder() {
        GetRomsFolder().then(result => {
            romsFolder = result
            getRoms(romsFolder)
        })
    }

    function getRoms(folder) {
        roms = []
        GetRoms(folder).then(result => {
            if (result != null) {
                for (let i = 0; i < result.length; i++) {
                    console.log(result[i])
                    roms = roms.concat({
                        value: result[i]
                    })
                }
                selectedRom = roms[0].value
            }

            savedefaults()
        })
    }

    function play() {
        Play(romsFolder, "TODO").then(result => {
            WindowHide()
        })
    }

    function savedefaults() {
        SaveDefaults(romsFolder, selectedRom, selectedMode)
    }

    EventsOn("EMULATOR_CLOSED", function() {
        WindowShow()
    })

    GetSettings().then(result => {
        romsFolder = result.romsFolder;
        getRoms(romsFolder)

        selectedRom = result.defaultRom
        selectedMode = result.defaultMode
    });
</script>

<main>
    <div class="window">
        <div class="logoContainer">
            <img alt="Iron logo" class="logo" src="{ironlogo}">
            <img alt="Launcher logo" class="logo" src="{launcherlogo}">
            <div class="padding"></div>
        </div>

        <div class="padding"></div>
        <div class="optionsContainer">
            <div>
                <div id="romsSection"> 
                    <div class="flex items-center space-x-2">
                        <Label for="romsSection" >Game:</Label>
                        <select bind:value={selectedRom} on:change={savedefaults}>
                            <option disabled value="-- GAME --">-- GAME --</option>
                            {#each roms as rom}
                                <option value={rom.value}>{rom.value}</option>
                            {/each}
                        </select>
                        <button class="btn" on:click={getRomsFolder}>ROMs</button>
                    </div>
                </div>
                <div class="padding"></div>
                <div id="modeSection">
                    <Label for="modeSection" >Mode:</Label>
                    <select bind:value={selectedMode} on:change={savedefaults} >
                        {#each modes as mode}
                            <option value={mode.value}>{mode.value}</option>
                        {/each}
                    </select>
                </div>
            </div>
            <div id="settingsSection">
                <Label for="settingsSection" >Settings</Label>
                {#each settings as setting}
                    <div class="flex items-center space-x-2">
                        <Checkbox id="settingSelect" />
                        <Label for="settingSelect" ><div contenteditable="true" bind:innerText={setting.value}></div></Label>
                    </div>
                {/each}
            </div>
        </div>
        <div class="padding"></div>
        <div class="play">
            <button class="btn" on:click={play}>Play</button>
        </div>
    </div>

</main>
  
<style>

    .window {
        background-color: rgb(138, 181, 216);
        width: 100%;
        height: 100%;
    }

    .play {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
    }

    .play > * {
        margin: auto;
        width: 100%;
    }

    .padding {
        margin-bottom: 2%;
    }

    .optionsContainer {
        height: 5rem;
        display: flex;
        align-items: flex-start;
        justify-content: space-evenly;
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