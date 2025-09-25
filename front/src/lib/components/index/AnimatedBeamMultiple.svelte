<script lang="ts">
  import { cn } from "$lib/utils";
  import { Github, IdCard, Key, User } from "@lucide/svelte";
  import AnimatedBeam from "./AnimatedBeam.svelte";
  import Circle from "./Circle.svelte";
  let containerRef = $state();
  let div1Ref = $state();
  let div2Ref = $state();
  let div3Ref = $state();
  let div4Ref = $state();
  let div5Ref = $state();
  let div6Ref = $state();
  let div7Ref = $state();
  let className: any = $state("");
  export { className as class };
</script>

<div
  bind:this={containerRef}
  class={cn("relative flex w-full items-center justify-center ", className)}
>
  <div
    class="flex h-full w-full flex-row justify-between gap-10 max-w-lg items-center"
  >
    <div class="flex flex-col justify-center gap-2">
      <!-- Div 1 -->
      <Circle bind:ref={div1Ref}>
        <div class="tooltip" data-tip="Leak of user's informations">
          <User />
        </div>
      </Circle>
      <!-- Div 2 -->
      <Circle bind:ref={div2Ref}>
        <div class="tooltip" data-tip="Leak of user's passwords">
          <Key />
        </div>
      </Circle>
      <!-- Div 3 -->
      <Circle bind:ref={div3Ref}>
        <div class="tooltip" data-tip="Leak of personal informations">
          <IdCard />
        </div>
      </Circle>
      <!-- Div 4 -->
      <Circle bind:ref={div4Ref}>
        <div class="tooltip" data-tip="Github recon">
          <Github />
        </div>
      </Circle>
      <!-- Div 5 -->
      <Circle bind:ref={div5Ref}>
        <div class="tooltip" data-tip="Google hunt">
          <svg
            width="16"
            viewBox="0 0 256 262"
            xmlns="http://www.w3.org/2000/svg"
            preserveAspectRatio="xMidYMid"
            class="fill-base-content"
            ><path
              d="M255.878 133.451c0-10.734-.871-18.567-2.756-26.69H130.55v48.448h71.947c-1.45 12.04-9.283 30.172-26.69 42.356l-.244 1.622 38.755 30.023 2.685.268c24.659-22.774 38.875-56.282 38.875-96.027"
            /><path
              d="M130.55 261.1c35.248 0 64.839-11.605 86.453-31.622l-41.196-31.913c-11.024 7.688-25.82 13.055-45.257 13.055-34.523 0-63.824-22.773-74.269-54.25l-1.531.13-40.298 31.187-.527 1.465C35.393 231.798 79.49 261.1 130.55 261.1"
            /><path
              d="M56.281 156.37c-2.756-8.123-4.351-16.827-4.351-25.82 0-8.994 1.595-17.697 4.206-25.82l-.073-1.73L15.26 71.312l-1.335.635C5.077 89.644 0 109.517 0 130.55s5.077 40.905 13.925 58.602l42.356-32.782"
            /><path
              d="M130.55 50.479c24.514 0 41.05 10.589 50.479 19.438l36.844-35.974C195.245 12.91 165.798 0 130.55 0 79.49 0 35.393 29.301 13.925 71.947l42.211 32.783c10.59-31.477 39.891-54.251 74.414-54.251"
            /></svg
          >
        </div>
      </Circle>
    </div>
    <div class="flex flex-col justify-center">
      <!-- Div 6 -->
      <Circle bind:ref={div6Ref}>
        <div class="tooltip" data-tip="Your eleakxir backend">
        <svg
          width={24}
          viewBox="0 0 141 205"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          class={cn("fill-primary", className)}
        >
          <path
            d="M69.7444 0C84.49 25.1637 100.559 49.4708 117.95 72.9219C126.451 85.0149 133.113 98.1035 137.934 112.188C144.168 135.735 140.195 157.355 126.014 177.046C109.938 196.591 89.1945 205.53 63.7844 203.865C36.717 200.867 17.4935 187.019 6.11353 162.321C-1.9191 142.522 -2.03583 122.655 5.76294 102.722C9.71019 93.7604 14.3849 85.2295 19.7864 77.1289C35.0593 56.0531 49.5504 34.4338 63.259 12.2705C65.6086 8.27249 67.7699 4.18207 69.7444 0ZM100.596 81.3359C102.957 92.649 102.198 103.751 98.3176 114.642C93.9276 124.99 87.7338 134.105 79.7366 141.987C77.6951 144.434 75.8254 147.005 74.1272 149.7C70.5033 155.43 68.5745 161.682 68.342 168.456C68.1692 175.079 70.6236 180.455 75.7043 184.583C89.1062 183.345 100.267 177.678 109.186 167.58C123.101 149.518 125.672 129.885 116.899 108.682C112.25 99.0223 106.815 89.9071 100.596 81.3359ZM70.095 36.4609C59.5617 53.4793 48.4018 70.0738 36.6145 86.2441C30.2619 94.8335 25.2366 104.183 21.5393 114.291C15.2642 135.159 19.2377 153.74 33.4592 170.034C39.7381 176.533 47.2756 180.798 56.0715 182.83C55.7506 182.699 55.4583 182.524 55.1956 182.305C40.6021 163.32 39.6671 143.687 52.3909 123.406C58.721 114.622 64.9149 105.74 70.9719 96.7617C72.559 94.0549 73.9614 91.2501 75.179 88.3477C81.8825 70.1929 80.1876 52.8971 70.095 36.4609Z"
          />
        </svg>
        </div>
      </Circle>
    </div>
    <div class="flex flec-col justify-center">
      <!-- Div 7 -->
      <Circle bind:ref={div7Ref}>
        <div class="tooltip" data-tip="This web client">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="lucide lucide-monitor-icon lucide-monitor"
          ><rect width="20" height="14" x="2" y="3" rx="2" /><line
            x1="8"
            x2="16"
            y1="21"
            y2="21"
          /><line x1="12" x2="12" y1="17" y2="21" /></svg
        >
        </div>
      </Circle>
    </div>
  </div>
  <AnimatedBeam bind:containerRef bind:fromRef={div1Ref} bind:toRef={div6Ref} />
  <AnimatedBeam bind:containerRef bind:fromRef={div2Ref} bind:toRef={div6Ref} />
  <AnimatedBeam bind:containerRef bind:fromRef={div3Ref} bind:toRef={div6Ref} />
  <AnimatedBeam bind:containerRef bind:fromRef={div4Ref} bind:toRef={div6Ref} />
  <AnimatedBeam bind:containerRef bind:fromRef={div5Ref} bind:toRef={div6Ref} />
  <AnimatedBeam bind:containerRef bind:fromRef={div6Ref} bind:toRef={div7Ref} />
</div>
