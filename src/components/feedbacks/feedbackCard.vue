<template functional>
  <b-card no-body class="feedback-card">
    <header>
      <h5 class="text-capitalize">
        <p class="m-0 feedback-title">{{props.feedback.title}}</p>
        <p class="text-muted m-0 feedback-date">{{props.feedback.update_at | date}}</p>
      </h5>
      <div
        class="text-muted text-truncate"
      >From {{props.feedback.display_name}} | {{props.feedback.creator}}</div>
    </header>
    <div class="feedback-body secondary-font">{{props.feedback.body}}</div>
  </b-card>
</template>
<style lang="scss">
.feedback-card {
  box-shadow: 0 1px 6px -2px rgba(32, 33, 36, 0.103);
  border-radius: 3px;

  header {
    margin: 0;
    display: flex;
    flex-direction: column;
    border-bottom: 1px solid #dee2e6;
    padding: 0.7rem 0.5rem 0.3rem;
    align-items: flex-start;
    background: #f8f9fa;

    h5 {
      display: flex;
      width: 100%;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 0.2rem;

      .feedback-title {
        font-size: 1.25rem;
      }

      .feedback-date {
        width: fit-content;
        white-space: nowrap;
        padding-left: 0.5rem;
      }
    }
  }
  .feedback-body {
    padding: 1rem 0.5rem;
    max-height: 150px;
    overflow: auto;
    min-height: 80px;
  }
}
</style>