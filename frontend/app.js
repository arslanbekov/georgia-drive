new Vue({
  el: "#app",
  data: {
    loading: true,
    lastExecutionTime: "",
    countdown: "",
    activeRoute: "/api/get-auto",
    activeTab: 1,
    selectedCityTab: 0,
    results: [],
    citiesOrder: [
      { name: "Kutaisi", order: 1 },
      { name: "Batumi", order: 2 },
      { name: "Telavi", order: 3 },
      { name: "Akhaltsikhe", order: 4 },
      { name: "Zugdidi", order: 5 },
      { name: "Gori", order: 6 },
      { name: "Poti", order: 7 },
      { name: "Ozurgeti", order: 8 },
      { name: "Sachkhere", order: 9 },
      { name: "Rustavi", order: 10 },
    ],
    currentRoute: "/api/get-auto",
    currentTitle: "Drive city (auto)",
  },
  created() {
    this.fetchData();
    this.fetchData();
    this.fetchLastExecutionTime();
  },
  methods: {
    cityResults(cityName) {
      return this.results.filter((result) => result.name === cityName);
    },
    isSaturday(dateString) {
      if (typeof dateString !== "string") return false;
      const date = new Date(dateString);
      return date.getDay() === 6;
    },
    setCurrentRoute(route) {
      this.loading = true;
      this.activeRoute = route;
      this.currentRoute = route;

      const activeTab = document.querySelector(".nav-tabs .nav-link.active");
      if (activeTab) {
        const parentLi = activeTab.closest("li");
        if (parentLi) {
          this.selectedCityTab = Array.from(
            parentLi.parentElement.children
          ).indexOf(parentLi);
        }
      }

      this.activeRoute = route;
      this.currentRoute = route;

      if (route === "/api/get-auto") {
        this.currentTitle = "Drive city (auto)";
      } else if (route === "/api/get-manual") {
        this.currentTitle = "Drive city (manual)";
      } else if (route === "/api/get-theory") {
        this.currentTitle = "Theory";
      }

      this.fetchData();

      this.$nextTick(() => {
        const tabLinks = this.$el.querySelectorAll(".nav-link");
        if (this.selectedCityTab !== null && tabLinks[this.selectedCityTab]) {
          new bootstrap.Tab(tabLinks[this.selectedCityTab]).show();
        } else {
          new bootstrap.Tab(tabLinks[0]).show();
        }
      });
    },
    fetchLastExecutionTime() {
      axios
        .get("/api/last-exec-time")
        .then((response) => {
          const timestamp = response.data.Timestamp;
          const dateObj = new Date(timestamp);
          this.lastExecutionTime = this.formatDate(dateObj);
          this.startCountdown(dateObj);
        })
        .catch((error) => {
          console.error("Error fetching last execution time:", error);
        });
    },

    startCountdown(startTime) {
      const endTime = new Date(startTime.getTime() + 12 * 60 * 1000);
      const updateCountdown = () => {
        const now = new Date();
        const difference = endTime - now;

        if (difference <= 0) {
          clearInterval(interval);
          this.countdown = "00:00";
        } else {
          const minutes = Math.floor(difference / (60 * 1000));
          const seconds = Math.floor((difference % (60 * 1000)) / 1000);

          this.countdown = `${String(minutes).padStart(2, "0")}:${String(
            seconds
          ).padStart(2, "0")}`;
        }
      };

      const interval = setInterval(updateCountdown, 1000);
      updateCountdown(); // вызываем сразу, чтобы обновить значение перед началом отсчета
    },
    formatDate(dateObj) {
      return `${dateObj.toLocaleDateString()} ${dateObj.toLocaleTimeString()}`;
    },
    fetchData() {
      fetch(this.currentRoute)
        .then((response) => response.json())
        .then((data) => {
          this.results = data;
          this.loading = false;

          this.$nextTick(() => {
            const tabLinks = this.$el.querySelectorAll(".nav-link");
            if (
              this.selectedCityTab !== null &&
              tabLinks[this.selectedCityTab]
            ) {
              new bootstrap.Tab(tabLinks[this.selectedCityTab]).show();
            } else {
              new bootstrap.Tab(tabLinks[0]).show();
            }
          });
        })
        .catch((error) => {
          console.error("There was an error fetching the data:", error);
          this.loading = false;
        });
    },
  },
});
