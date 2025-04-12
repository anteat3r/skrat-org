<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import Graph from './lib/components/Graph.svelte';
	import Knob from './lib/components/Knob.svelte';
	import BallBalancer from './lib/components/BallBalancer.svelte';
	import Toggle from './lib/components/Toggle.svelte';

	// --- State ---
	let currentMode: string = 'normal_mode';
	let isOnline: boolean = true; // Simulate connection status

	let isControllingFromHere: boolean = false;

	// PID constants state: These reflect user input when controlling, or received data otherwise
	let pidConstants = { p: 12.8, i: 0.007, d: 28.48 };
	// Stores the latest values received from the source
	let receivedPidConstants = { p: 12.8, i: 0.007, d: 28.48 };
	// Tracks which PID terms are enabled/disabled locally
	let pidEnabled = { p: true, i: true, d: false }; // Initial enabled state matches image

	// Data for graphs
	let errorData: number[] = [];
	let controlInputData: number[] = [];
	const MAX_GRAPH_POINTS = 100; // Max points to display on graphs

	// Current values for display
	let currentError: number | null = null;
	let currentControlInput: number | null = null;

	// Ball balancer state
	let ballPosition: number = 0; // -1 (left) to 1 (right)
    let leverAngle: number = 0; // Degrees for visualization tilt

	// --- Simulation ---
	let simulationInterval: number | undefined = undefined;

  async function connect() {
    try {
    let device = await (navigator as any).bluetooth.requestDevice({filters: [{services: ["6e400001-b5a3-f393-e0a9-e50e24dcca9e"]}]});
      await device.gatt.connect();
      let service = await device.gatt.getPrimaryService("6e400001-b5a3-f393-e0a9-e50e24dcca9e");
      let characteristic = await service.getCharacteristic("6e400002-b5a3-f393-e0a9-e50e24dcca9e");
      let idk = await characteristic.startNotifications();
      let decoder = new TextDecoder();
      idk.addEventListener('characteristicvaluechanged', function (evt: { target: { value: any; }; }) {
        let value = evt.target.value;
        let str = decoder.decode(value);
        let ss = str.split("|");
        switch (ss[0]) {
          case "dist":
            let dist = parseInt(ss[1]);
            ballPosition = dist / 300 - 1;
            currentControlInput = ballPosition;
            controlInputData.push(ballPosition);
          break;
          case "error":
            let error = parseInt(ss[1]);
            currentError = error;
            errorData.push(error);
          break;
          case "pwm":
            let pwm = parseInt(ss[1]);
            leverAngle = pwm - 51;
          break;
          case "knob1":
            let knob1 = parseInt(ss[1]);
            pidConstants.p = knob1 / 1024 * 50;
          break;
          case "knob2":
            let knob2 = parseInt(ss[1]);
            pidConstants.i = knob2 / 1024 * .1;
          break;
          case "knob3":
            let knob3 = parseInt(ss[1]);
            pidConstants.d = knob3 / 1024 * 100;
          break;
        }
      })
    } catch (e) {
      alert(e);
    }
  }

	// --- Communication Placeholder ---
	/**
	 * Placeholder function to simulate sending PID data.
	 * Replace this with your actual communication logic (e.g., Web Bluetooth, Web Serial, HTTP POST).
	 */
	function sendPidUpdate(p: number, i: number, d: number, pEnabled: boolean, iEnabled: boolean, dEnabled: boolean) {
		// The controller (microbit) should use 0 if the term is disabled.
        const effectiveP = pEnabled ? p : 0;
        const effectiveI = iEnabled ? i : 0;
        const effectiveD = dEnabled ? d : 0;

        console.log('Sending PID Update:', {
			p: effectiveP,
			i: effectiveI,
			d: effectiveD,
            p_raw: p, // Send raw too maybe?
            i_raw: i,
            d_raw: d,
            p_enabled: pEnabled,
            i_enabled: iEnabled,
            d_enabled: dEnabled
		});
		// --- Add your actual sending code here ---
		// Example: await sendDataToMicrobit({ type: 'pid_update', p: effectiveP, i: effectiveI, d: effectiveD });
	}

	// --- Event Handlers ---
	function handleControlToggle(event: CustomEvent<boolean>) {
		isControllingFromHere = event.detail;
        if(isControllingFromHere) {
            // When switching TO control, ensure knobs start from last known received value
            pidConstants = { ...receivedPidConstants };
        } else {
            // When switching OFF control, immediately reflect the received values ("teleport")
             pidConstants = { ...receivedPidConstants };
             // Optional: If the remote also sends enable status, sync here too
             // pidEnabled = { ...receivedPidEnabledStatus };
        }
	}

    function handlePidEnableToggle(term: 'p' | 'i' | 'd', event: CustomEvent<boolean>) {
        pidEnabled[term] = event.detail;
         if (isControllingFromHere) {
            // Send update immediately if controlling from here and enable status changes
            sendPidUpdate(pidConstants.p, pidConstants.i, pidConstants.d, pidEnabled.p, pidEnabled.i, pidEnabled.d);
        }
    }

	function handleKnobChange(term: 'p' | 'i' | 'd', event: CustomEvent<number>) {
        // This check prevents feedback loops if the knob updates reactively
        if (!isControllingFromHere) return;

		pidConstants = { ...pidConstants, [term]: event.detail };
		// Send the update to the microbit/backend
		sendPidUpdate(pidConstants.p, pidConstants.i, pidConstants.d, pidEnabled.p, pidEnabled.i, pidEnabled.d);
	}

	// --- Lifecycle ---
	onMount(() => {
    connect();
	});

	onDestroy(() => {
		if (simulationInterval) {
			clearInterval(simulationInterval);
		}
	});

</script>

<div class="app-container">
	<header>
		<h1>PID control</h1>
		<div class="status">
			<span>{currentMode}</span>
			<span class="status-dot" class:online={isOnline}></span>
		</div>
	</header>

	<main class="content-grid">
		<!-- Left Column: Controls -->
		<section class="controls-column">
			<!-- <div class="control-toggle-section"> -->
			<!-- 	<Toggle bind:checked={isControllingFromHere} on:change={handleControlToggle} /> -->
			<!-- 	<span>Control from here</span> -->
			<!-- </div> -->

            <!-- Proportional -->
            <div class="pid-control-block">
                 <div class="pid-header">
                     <Toggle
                        bind:checked={pidEnabled.p}
                        disabled={!isControllingFromHere}
                        on:change={(e) => handlePidEnableToggle('p', e)}
                     />
                     <span>Proportional</span>
                 </div>
                 <Knob
                    label="Proportional"
                    bind:value={pidConstants.p}
                    min={0} max={50} step={0.1}
                    readOnly={!isControllingFromHere}
                    disabled={!pidEnabled.p && isControllingFromHere}
                    on:change={(e: CustomEvent<number>) => handleKnobChange('p', e)}
                 />
            </div>

             <!-- Integral -->
             <div class="pid-control-block">
                 <div class="pid-header">
                     <Toggle
                        bind:checked={pidEnabled.i}
                        disabled={!isControllingFromHere}
                        on:change={(e) => handlePidEnableToggle('i', e)}
                    />
                     <span>Integral</span>
                 </div>
                 <Knob
                    label="Integral"
                    bind:value={pidConstants.i}
                    min={0} max={0.1} step={0.001}
                    readOnly={!isControllingFromHere}
                     disabled={!pidEnabled.i && isControllingFromHere}
                    on:change={(e) => handleKnobChange('i', e)}
                 />
             </div>

            <!-- Derivative -->
             <div class="pid-control-block">
                 <div class="pid-header">
                     <Toggle
                        bind:checked={pidEnabled.d}
                        disabled={!isControllingFromHere}
                        on:change={(e) => handlePidEnableToggle('d', e)}
                    />
                     <span>Derivative</span>
                 </div>
                 <Knob
                    label="Derivative"
                    bind:value={pidConstants.d}
                    min={0} max={100} step={0.01}
                    readOnly={!isControllingFromHere}
                    disabled={!pidEnabled.d && isControllingFromHere}
                    on:change={(e) => handleKnobChange('d', e)}
                 />
             </div>

		</section>

		<!-- Right Column: Visualizations -->
		<section class="visuals-column">
			<BallBalancer {ballPosition} leverAngleDeg={leverAngle} />

			<Graph
				label="Error"
				unit="cm"
				color="#dc3545"
                fill={true}
                targetValue={0}
				data={errorData}
				currentValue={currentError}
                maxPoints={MAX_GRAPH_POINTS}
			/>

            <!-- Placeholder for connection arrow/icon -->
            <div class="connection-arrow">
                <!-- You can use an SVG or an icon font here -->
                 🧠 ↓
            </div>

			<Graph
				label="Control input"
				unit="deg"
				color="#28a745"
                fill={true}
                targetValue={0}
				data={controlInputData}
				currentValue={currentControlInput}
                maxPoints={MAX_GRAPH_POINTS}
			/>
		</section>
	</main>
</div>

<style>
	/* Basic App Structure & Styling */
	:global(body) {
		margin: 0;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
		background-color: #f0f2f5; /* Light background */
        color: #333;
	}

	.app-container {
		max-width: 1000px; /* Max width */
		margin: 0 auto; /* Center */
		padding: 20px;
        background-color: #ffffff; /* White container background */
        min-height: 100vh;
        display: flex;
        flex-direction: column;
        box-shadow: 0 2px 10px rgba(0,0,0,0.1);
	}

	header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		border-bottom: 1px solid #e0e0e0;
		padding-bottom: 15px;
		margin-bottom: 20px;
	}

	header h1 {
		margin: 0;
        font-size: 1.8em;
        font-weight: 500;
	}

	.status {
		display: flex;
		align-items: center;
		font-size: 0.9em;
        color: #555;
	}

	.status-dot {
		width: 10px;
		height: 10px;
		background-color: #dc3545; /* Red (offline) */
		border-radius: 50%;
		margin-left: 8px;
		display: inline-block;
	}

	.status-dot.online {
		background-color: #28a745; /* Green (online) */
	}

    .content-grid {
        display: grid;
        grid-template-columns: 200px 1fr; /* Fixed width for controls, rest for visuals */
        gap: 30px; /* Gap between columns */
        flex-grow: 1;
    }

    .controls-column {
       display: flex;
       flex-direction: column;
       gap: 15px; /* Space between control blocks */
    }

     .visuals-column {
       display: flex;
       flex-direction: column;
       align-items: center; /* Center items like arrow */
    }

    .control-toggle-section {
        display: flex;
        align-items: center;
        gap: 10px;
        margin-bottom: 20px;
        padding: 10px;
        background-color: #f8f9fa;
        border-radius: 8px;
    }

    .pid-control-block {
        background-color: #f8f9fa;
        padding: 15px;
        border-radius: 8px;
        box-shadow: 0 1px 2px rgba(0,0,0,0.05);
    }

    .pid-header {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 10px;
        font-weight: 500;
    }

    .connection-arrow {
        font-size: 2em;
        color: #888;
        margin: -10px 0; /* Adjust vertical spacing */
        text-align: center;
    }

    /* Responsive adjustments if needed */
    @media (max-width: 768px) {
        .content-grid {
            grid-template-columns: 1fr; /* Stack columns on smaller screens */
        }
        .controls-column {
            flex-direction: row; /* Lay out controls horizontally */
            flex-wrap: wrap;
            justify-content: center;
        }
        .visuals-column {
            margin-top: 20px;
        }
        .app-container {
            padding: 10px;
        }
        header h1 {
            font-size: 1.5em;
        }
    }

</style>
