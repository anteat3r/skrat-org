<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	export let ballPosition: number = 0; // Normalized position, e.g., -1 (left) to 1 (right)
	export let leverAngleDeg: number = 0; // Angle of the lever

	let canvas: HTMLCanvasElement;
	let ctx: CanvasRenderingContext2D | null = null;
	let width: number = 300;
	let height: number = 100;
    let container: HTMLDivElement;

    // --- Constants for drawing ---
    const leverHeight = 8;
    const ballRadius = 12; // Slightly smaller maybe?
    const fulcrumWidth = 30;
    const fulcrumHeight = 20;

	function drawBalancer() {
		if (!ctx || !canvas || !width || !height) return; // Ensure context and dimensions are ready

        const leverWidth = width * 0.8;
		const centerX = width / 2;
		const centerY = height * 0.6; // Position fulcrum lower down

		// Clear canvas
		ctx.clearRect(0, 0, width, height);

        // --- Draw Fulcrum (Triangle) - Drawn first, not rotated ---
        ctx.fillStyle = '#000';
        ctx.beginPath();
        ctx.moveTo(centerX - fulcrumWidth / 2, centerY + leverHeight / 2); // Start below estimated lever center
        ctx.lineTo(centerX + fulcrumWidth / 2, centerY + leverHeight / 2);
        ctx.lineTo(centerX, centerY + leverHeight / 2 + fulcrumHeight);
        ctx.closePath();
        ctx.fill();


		// --- Draw Lever and Ball (Rotated) ---
        ctx.save(); // Save context state (transformations, styles)

        // Apply transformations: Move origin to pivot, then rotate
		ctx.translate(centerX, centerY);
		ctx.rotate(leverAngleDeg * Math.PI / 180);

		// Draw Lever centered horizontally around the new (0,0) origin
		ctx.fillStyle = '#333'; // Darker lever
		ctx.fillRect(-leverWidth / 2, -leverHeight / 2, leverWidth, leverHeight);

		// Calculate Ball Position relative to the *rotated* lever's center (0,0)
		// Map ballPosition (-1 to 1) to the available space on the lever
        const maxBallTravel = leverWidth / 2 - ballRadius; // Max distance from center
		const ballX_local = ballPosition * maxBallTravel; // X pos in rotated frame
		const ballY_local = -leverHeight / 2 - ballRadius; // Y pos in rotated frame (on top of lever)

		// Draw Ball at calculated local position
		ctx.fillStyle = '#3498db'; // Brighter blue ball
		ctx.strokeStyle = '#2980b9'; // Slightly darker blue outline
		ctx.lineWidth = 2;
		ctx.beginPath();
		// Use arc(x, y, radius, startAngle, endAngle)
		ctx.arc(ballX_local, ballY_local, ballRadius, 0, Math.PI * 2);
		ctx.fill();
		ctx.stroke();

		ctx.restore(); // Restore context state (removes translate and rotate)
	}

    // --- Resize Handling ---
    let resizeObserver: ResizeObserver;
    function setupResizeObserver() {
        if (!container || typeof ResizeObserver === 'undefined') return;
        resizeObserver = new ResizeObserver(entries => {
            // Debounce or throttle resize updates if performance becomes an issue
            for (let entry of entries) {
                const rect = entry.contentRect;
                // Check if size actually changed to avoid unnecessary redraws
                if (rect.width !== width || rect.height !== height) {
                    width = rect.width;
                    height = rect.height;
                    canvas.width = width; // Update canvas resolution
                    canvas.height = height;
                    // No need to call drawBalancer here, reactive update will catch it
                }
            }
        });
        resizeObserver.observe(container);
    }

	onMount(() => {
        if(canvas){
		    ctx = canvas.getContext('2d');
            // Initial sizing from container
            width = container.clientWidth;
            height = container.clientHeight; // Use container height (respects CSS)
            canvas.width = width;
            canvas.height = height;
            setupResizeObserver(); // Setup observer after initial size
		    // drawBalancer(); // Initial draw will happen via reactive statement
        }
	});

     onDestroy(() => {
       if (resizeObserver && container) {
           resizeObserver.unobserve(container);
       }
    });

    // --- Reactive Redraw ---
    // IMPORTANT: Make this depend on the props that affect the drawing!
	$: if (ctx && width && height && ballPosition !== undefined && leverAngleDeg !== undefined) {
        // console.log('Redrawing balancer:', { ballPosition, leverAngleDeg }); // For debugging
        requestAnimationFrame(drawBalancer); // Use rAF for smoother animation
    }

</script>

<div bind:this={container} class="balancer-container">
	<canvas bind:this={canvas}></canvas>
</div>

<style>
	.balancer-container {
		width: 100%;
		height: 150px; /* Adjust height as needed */
        margin-bottom: 20px;
	}
    canvas {
        display: block;
        width: 100%;
        height: 100%;
    }
</style>