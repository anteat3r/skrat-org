<div class="container">
    <div id="logo-wrapper">
        <img src="/images/scrat_logo_cnb.png" alt="logo" id="logo"/>
        <canvas id="logo-text-canvas" height="50" width="70"></canvas>
    </div>
    <div id="content-wrapper">
        <div id="about-us" class="nav-element">
            <a href="/about-us">About us</a>
            <i class="fa-solid fa-user nav-icon"></i>

        </div>
    </div>
</div>

<style>
    .container {
        box-sizing: border-box;
        background-color: #090306;
        width: 100%;
        height: auto;
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        height: 70px;
        padding-inline: 10px;
    }
    #logo-wrapper{
        display: flex;
        flex-direction: row;
        align-items: center;
        cursor: pointer;

        gap: 3px;
    }

    #logo{
        max-width: 100px;
        max-height: 60px;
        object-fit: contain;

    }

    #content-wrapper{
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: space-around;
        flex-grow: 1;
    }


    .nav-element{
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: space-between;
        cursor: pointer;
        gap: 6px;
    }

    .nav-element a{
        color: #e0e0e0;
        font-size: 17px;
        text-decoration: none;
    }

    .nav-icon{
        font-size: 15px;
    }

    #about-us a{
        color: #e0e0e0;
        text-decoration: none;
        transition: all 0.15s ease-out;
    }
    #about-us:hover a{
        text-decoration: underline;
        text-decoration-thickness: 2px;
        text-underline-offset: 3px;
        color: #b5c5d4;
        transition: all 0.15s ease-out;
    }
    
    #about-us i{
        color: #9cbcda;
        transition: all 0.15s ease-out;
    }
    #about-us:hover i{
        color: #3386d4;
        transition: all 0.15s ease-out;
    }




</style>

<script lang="ts">
    // setup the context
    let logo_canvas: HTMLCanvasElement | null;
    let logo_ctx: CanvasRenderingContext2D | null;

    // setup the letters that will be drawn
    let letters = ["S", "K", "R", "A", "T"];
    let letter_xoftsets = new Array(letters.length).fill(0);
    let letter_widths = new Array(letters.length).fill(0);

    let target_xoftset = 0;

    const text_oftset_max = 7;
    const text_opacity_min = 0.4;

    // bind loding of the page with setting up the canvas
    document.addEventListener('DOMContentLoaded', () => {
        logo_canvas = document.getElementById('logo-text-canvas') as HTMLCanvasElement | null;
        if (logo_canvas) {
            logo_ctx = logo_canvas.getContext('2d');
            if (logo_ctx) {
                // set the font's style
                logo_ctx.font = 'bold 17px Lexend';
                logo_ctx.textAlign = "left";
                logo_ctx.textBaseline = "middle";
                logo_ctx.fillStyle = "#e0e0e0";

                // measure the letters length
                for (let i = 0; i < letters.length; i++) {
                    letter_widths[i] = logo_ctx.measureText(letters[i]).width;
                }
                
                // call the first frame
                update();
            } else {
                console.error("Unable to get 2D context for the canvas.");
            }
        } else {
            console.error("Canvas element with ID 'logo-text-canvas' not found.");
        }
        
        document.getElementById('logo-wrapper')!.addEventListener('mouseover', () => target_xoftset = text_oftset_max);
        document.getElementById('logo-wrapper')!.addEventListener('mouseout', () => target_xoftset = 0);
    });

    function update() {
        requestAnimationFrame(update);

        animation_update();

        if (logo_ctx && logo_canvas) {
            logo_ctx.clearRect(0, 0, logo_canvas.width, logo_canvas.height);
            let center_x = logo_canvas.width / 2;
            let center_y = logo_canvas.height / 2;
            let total_width = 0;
            for (let i = 0; i < letters.length; i++) {
                total_width += letter_widths[i];
            }
            let start_x = 0;
            for (let i = 0; i < letters.length; i++) {
                logo_ctx.fillStyle = `rgba(255, 255, 255, ${Math.max(0, Math.min(1, (letter_xoftsets[i]) / (text_oftset_max)*(1 - text_opacity_min) + text_opacity_min))})`;
                
                logo_ctx.fillText(letters[i], start_x + letter_xoftsets[i], center_y);
                start_x += letter_widths[i];
            }
        }
    }
    
    // updates the position of the letters
    // (their x oftset) based on the target x oftset which is
    // set by weather the user is hovering over the logo wrapper
    function animation_update() {
        for (let i = 0; i < letters.length; i++) {
            if (i === 0) {
                letter_xoftsets[0] += (target_xoftset - letter_xoftsets[0]) * 0.2;
            } else {
                letter_xoftsets[i] += (letter_xoftsets[i - 1] - letter_xoftsets[i]) * 0.4;
            }
        }
    }






</script>