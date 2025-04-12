<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let checked: boolean = false;
	export let disabled: boolean = false;

	const dispatch = createEventDispatcher<{ change: boolean }>();

	function handleChange(event: Event) {
		const target = event.target as HTMLInputElement;
		checked = target.checked;
		dispatch('change', checked);
	}
</script>

<label class="toggle-switch" class:disabled>
	<input type="checkbox" {checked} {disabled} on:change={handleChange} />
	<span class="slider" />
</label>

<style>
	/* Basic Toggle Switch Styles (adapt as needed) */
	.toggle-switch {
		position: relative;
		display: inline-block;
		width: 40px; /* Smaller width */
		height: 20px; /* Smaller height */
		cursor: pointer;
	}
    .toggle-switch.disabled {
        cursor: not-allowed;
        opacity: 0.6;
    }

	.toggle-switch input {
		opacity: 0;
		width: 0;
		height: 0;
	}

	.slider {
		position: absolute;
		cursor: pointer;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: #ccc;
		transition: 0.4s;
		border-radius: 20px; /* Fully rounded */
	}

	.slider:before {
		position: absolute;
		content: '';
		height: 16px; /* Smaller handle */
		width: 16px; /* Smaller handle */
		left: 2px; /* Adjust position */
		bottom: 2px; /* Adjust position */
		background-color: white;
		transition: 0.4s;
		border-radius: 50%;
	}

	input:checked + .slider {
		background-color: #5cb85c; /* Green when checked */
	}
     input:disabled + .slider {
        background-color: #e0e0e0;
    }
     input:disabled:checked + .slider {
        background-color: #a5d6a7; /* Lighter green when disabled checked */
    }
     input:disabled + .slider:before {
        background-color: #f5f5f5;
    }


	input:checked + .slider:before {
		transform: translateX(20px); /* Adjust translation distance */
	}
</style>