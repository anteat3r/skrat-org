<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	export let data: number[] = [];
	export let label: string = 'Data';
	export let unit: string = '';
	export let currentValue: number | null = null;
	export let color: string = '#007bff'; // Default color
    export let targetValue: number = 0; // The dashed line target
    export let fill: boolean = true; // Whether to fill below the line
    export let maxPoints: number = 100; // Max data points to show

	let canvas: HTMLCanvasElement;
	let ctx: CanvasRenderingContext2D | null = null;
	let width: number = 300; // Initial default
	let height: number = 150; // Initial default
    let container: HTMLDivElement;

    let rafId: number;

	function drawGraph() {
		if (!ctx || !canvas) return;

		ctx.clearRect(0, 0, width, height);

        const padding = { top: 20, right: 50, bottom: 20, left: 10 };
        const plotWidth = width - padding.left - padding.right;
        const plotHeight = height - padding.top - padding.bottom;

        if (data.length < 2) return; // Need at least 2 points to draw a line

        // Determine data range dynamically
        const visibleData = data.slice(-maxPoints);
        let minVal = Math.min(...visibleData, targetValue);
        let maxVal = Math.max(...visibleData, targetValue);
        const range = maxVal - minVal;
        const dataSpan = range === 0 ? 1 : range; // Avoid division by zero
        // Add some vertical padding to the range
        minVal -= dataSpan * 0.1;
        maxVal += dataSpan * 0.1;
        const effectiveRange = maxVal - minVal || 1; // Ensure range is not zero


        // --- Draw Target Line ---
        const targetY = padding.top + plotHeight - ((targetValue - minVal) / effectiveRange) * plotHeight;
        if (targetY >= padding.top && targetY <= padding.top + plotHeight) {
             ctx.strokeStyle = '#aaa';
             ctx.lineWidth = 1;
             ctx.setLineDash([5, 5]);
             ctx.beginPath();
             ctx.moveTo(padding.left, targetY);
             ctx.lineTo(padding.left + plotWidth, targetY);
             ctx.stroke();
             ctx.setLineDash([]); // Reset line dash
        }


		// --- Draw Data Line ---
        ctx.strokeStyle = color;
        ctx.lineWidth = 2;
        ctx.beginPath();

        visibleData.forEach((point, index) => {
            const x = padding.left + (index / (maxPoints - 1)) * plotWidth;
			const y = padding.top + plotHeight - ((point - minVal) / effectiveRange) * plotHeight;

            if (index === 0) {
                ctx?.moveTo(x, y);
            } else {
                ctx?.lineTo(x, y);
            }
        });
        ctx.stroke();

        // --- Draw Fill ---
        if (fill && visibleData.length > 0) {
            ctx.lineTo(padding.left + plotWidth, padding.top + plotHeight); // Bottom right
            ctx.lineTo(padding.left, padding.top + plotHeight); // Bottom left
            const firstY = padding.top + plotHeight - ((visibleData[0] - minVal) / effectiveRange) * plotHeight;
            ctx.lineTo(padding.left, firstY); // Back to start Y
            ctx.closePath();
            const gradient = ctx.createLinearGradient(0, padding.top, 0, padding.top + plotHeight);
            // Adjust gradient color based on main color (simple example)
            const rgbColor = hexToRgb(color);
            if (rgbColor) {
                 gradient.addColorStop(0, `rgba(${rgbColor.r}, ${rgbColor.g}, ${rgbColor.b}, 0.3)`);
                 gradient.addColorStop(1, `rgba(${rgbColor.r}, ${rgbColor.g}, ${rgbColor.b}, 0.05)`);
            } else {
                 gradient.addColorStop(0, 'rgba(0,0,0,0.1)');
                 gradient.addColorStop(1, 'rgba(0,0,0,0.01)');
            }

            ctx.fillStyle = gradient;
            ctx.fill();
        }


        // --- Draw Current Value Indicator ---
        if (visibleData.length > 0) {
            const lastX = padding.left + plotWidth;
            const lastY = padding.top + plotHeight - ((visibleData[visibleData.length - 1] - minVal) / effectiveRange) * plotHeight;
            ctx.fillStyle = color;
            ctx.beginPath();
            ctx.arc(lastX, lastY, 4, 0, Math.PI * 2);
            ctx.fill();

            // Optionally draw a small circle outline
             ctx.strokeStyle = '#fff'; // White outline for contrast
             ctx.lineWidth = 1.5;
             ctx.beginPath();
             ctx.arc(lastX, lastY, 5, 0, Math.PI * 2);
             ctx.stroke();
        }


		// --- Draw Labels ---
		ctx.fillStyle = '#333';
		ctx.font = '12px Arial';
		ctx.textAlign = 'left';
		ctx.fillText(label, padding.left, padding.top - 5);

		if (currentValue !== null) {
			ctx.textAlign = 'right';
			ctx.fillText(`${currentValue.toFixed(unit === 'deg' ? 3 : (unit === 'cm' ? 1: 2))} ${unit}`, width - 5, padding.top - 5);
		}
	}

    // Helper to convert hex to RGB
    function hexToRgb(hex: string): {r: number, g: number, b: number} | null {
        const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
        return result ? {
            r: parseInt(result[1], 16),
            g: parseInt(result[2], 16),
            b: parseInt(result[3], 16)
        } : null;
    }


    // Resize observer
    let resizeObserver: ResizeObserver;
    function setupResizeObserver() {
        if (!container) return;
        resizeObserver = new ResizeObserver(entries => {
            for (let entry of entries) {
                const rect = entry.contentRect;
                width = rect.width;
                height = rect.height;
                canvas.width = width;
                canvas.height = height;
                // No need to call drawGraph here, reactive update will handle it
            }
        });
        resizeObserver.observe(container);
    }

	onMount(() => {
        if (canvas) {
		    ctx = canvas.getContext('2d');
            width = container.clientWidth;
            height = container.clientHeight;
            canvas.width = width;
            canvas.height = height;
            setupResizeObserver();
		    drawGraph(); // Initial draw
        }
        // Use rAF loop triggered by data changes if performance needed,
        // but reactive $: drawGraph() is often sufficient for moderate updates.
	});

    onDestroy(() => {
       if (resizeObserver && container) {
           resizeObserver.unobserve(container);
       }
       cancelAnimationFrame(rafId); // Clean up animation frame if used
    });

    // Reactive statement to redraw whenever data or dimensions change
    $: if (ctx && width && height && data) {
         // Throttle drawing with requestAnimationFrame if updates are very frequent
         cancelAnimationFrame(rafId);
         rafId = requestAnimationFrame(drawGraph);
         // Or simpler: drawGraph();
    }

</script>

<div bind:this={container} class="graph-container">
	<canvas bind:this={canvas}></canvas>
</div>

<style>
	.graph-container {
		width: 100%;
		height: 150px; /* Default height, can be overridden */
        background-color: #fff;
        border-radius: 8px;
        box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        padding: 10px; /* Add padding inside the container */
        box-sizing: border-box; /* Include padding in width/height */
        margin-bottom: 20px;
        position: relative; /* Needed for absolute positioning of elements inside if required */
	}
    canvas {
        display: block;
        width: 100%;
        height: 100%;
    }
</style>