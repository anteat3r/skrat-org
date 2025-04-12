<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let value: number = 0;
	export let min: number = 0;
	export let max: number = 100;
	export let step: number = 0.1;
	export let label: string = '';
    export let disabled: boolean = false; // Used to disable interaction
    export let readOnly: boolean = false; // Used when not controlling from here

	const dispatch = createEventDispatcher<{ change: number }>();

	function handleInput(event: Event) {
		const target = event.target as HTMLInputElement;
		value = parseFloat(target.value);
		dispatch('change', value);
	}

    // Reactive update if the value prop changes from the parent
    $: if (typeof window !== 'undefined') { // Avoid SSR errors
        const inputElement = document.getElementById(`knob-${label}`) as HTMLInputElement;
        if (inputElement && parseFloat(inputElement.value) !== value) {
            inputElement.value = String(value);
        }
    }

</script>

<div class="knob-container" class:disabled={disabled || readOnly}>
    <!-- Basic visual representation - replace with SVG/Canvas later for real knob -->
     <div class="knob-visual">
        <!-- Placeholder for visual knob - could be SVG -->
        <svg viewBox="0 0 100 100" width="50" height="50">
            <circle cx="50" cy="50" r="45" fill="none" stroke="#e0e0e0" stroke-width="10"/>
            <line
                x1="50"
                y1="50"
                x2="50"
                y2="10"
                stroke="#555"
                stroke-width="8"
                stroke-linecap="round"
                transform={`rotate(${((value - min) / (max - min)) * 270 - 135} 50 50)`}
            />
        </svg>
     </div>
	<input
        id="knob-{label}"
		type="range"
		{min}
		{max}
		{step}
		bind:value={value}
        disabled={disabled || readOnly}
		on:input={handleInput}
        aria-label={label}
	/>
	<div class="value-display">{value.toFixed(label === 'Derivative' ? 2 : (label === 'Integral' ? 3 : 1))}</div>
</div>

<style>
	.knob-container {
		display: flex;
        flex-direction: column;
		align-items: center;
        padding: 10px;
        border-radius: 8px;
        background-color: #fff;
        box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        width: 120px; /* Adjust as needed */
        margin-bottom: 15px;
	}

    .knob-container.disabled {
       opacity: 0.7;
    }

    .knob-visual {
        margin-bottom: 5px;
        /* Basic pointer imitation */
    }

    .knob-visual svg {
        display: block;
        margin: 0 auto;
    }

	input[type='range'] {
        /* Hide the default range input for now, rely on visual */
		/* appearance: none; */
		width: 80%;
		height: 5px;
		background: #ddd;
		outline: none;
		opacity: 0.7;
		transition: opacity 0.2s;
        margin-top: 10px;
        margin-bottom: 5px;
        cursor: pointer;
	}
    input[type='range']:disabled {
        cursor: not-allowed;
    }


	input[type='range']::-webkit-slider-thumb {
		appearance: none;
		width: 15px;
		height: 15px;
		background: #5cb85c;
		cursor: pointer;
        border-radius: 50%;
	}

	input[type='range']::-moz-range-thumb {
		width: 15px;
		height: 15px;
		background: #5cb85c;
		cursor: pointer;
        border-radius: 50%;
        border: none;
	}

    .value-display {
        margin-top: 5px;
        font-size: 0.9em;
        color: #555;
        background-color: #f8f9fa;
        padding: 2px 8px;
        border-radius: 4px;
        min-width: 50px;
        text-align: center;
    }
</style>