<template>
  <div>
    <div class="auth-background"></div>
    <div class="fullsize">
        <div class="auth-container">
            <div class="branding">
                <h2 class="title is-5 has-text-white">{{ $store.state.networkName }}</h2>
            </div>
            <router-view></router-view>
        </div>
    </div>
  </div>
</template>

<style lang="scss">
.auth-background {
  position: fixed;
  right: 0;
  bottom: 0;
  min-width: 100%;
  min-height: 100%;
  background: rgb(0, 144, 255);
}

.auth-container {
  max-width: 330px;
  width: 80%;
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
}

.branding {
  .platform {
    position: relative;
    margin-left: 50px;
    top: 24px;
    transform: translate(0%, -50%);
    text-align: center;
    color: #fff;
  }
}

.form-auth {
  padding: 30px;
  box-shadow: rgba(0, 0, 0, 0.2) 0px 0px 20px;
  border: none;
  border-radius: 3px;
  background: rgba(255, 255, 255, 1);
  clear: both;

  h2 {
    margin-bottom: 20px;
  }

  .input {
    border-radius: 0;
    margin-bottom: -1px;
  }

  .input:first-of-type {
    border-top-left-radius: 4px;
    border-top-right-radius: 4px;
  }

  .input:last-of-type {
    border-bottom-right-radius: 4px;
    border-bottom-left-radius: 4px;
    margin-bottom: 10px;
  }

  .input:focus {
    z-index: 3;
  }
  .input:hover {
    z-index: 2;
  }
}
</style>

