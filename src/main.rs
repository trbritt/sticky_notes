use gdk::Display;
use gtk::{gio, glib, prelude::*};
use gio::ActionEntry;
use gtk::{gdk, Application, ApplicationWindow, Button, CssProvider, Orientation};

use std::cell::Cell;
use std::rc::Rc;
use glib::clone;
use glib::closure_local;
use chrono::{Datelike, Local, Weekday, NaiveDate};

use std::process::Command;
use std::io::{self,Write};

fn show_modal_window(d: u32) { //use this to launch the sticky notes
    let output = Command::new("/usr/bin/gnome-terminal")
    .arg("--geometry")
    .arg("100x200")
    .arg("--")
    .arg("/home/trbritt/Desktop/sticky_notes/driver/gonotes_driver")
    .arg(format!("id={d}"))
    .output()
    .expect("Failed to execute stickynotes");

    println!("launching sticky id={}, status: {}", d,output.status);
    io::stdout().write_all(&output.stdout).unwrap();
    io::stderr().write_all(&output.stderr).unwrap();
    // let modal = gtk::Window::builder()
    // .title(format!("{d}"))
    // .modal(false)
    // .default_height(100)
    // .default_width(100)
    // .build();

    // // modal.set_parent(parent);
    // let content_area = gtk::Text::builder()
    // .text("This is a modal!")
    // .build();

    // modal.set_child(Some(&content_area));
    // modal.present();
}

#[inline]
const fn weekday_to_u32(n: Weekday) -> u32 {
    match n  {
        Weekday::Mon => 1,
        Weekday::Tue => 2,
        Weekday::Wed => 3,
        Weekday::Thu => 4,
        Weekday::Fri => 5,
        Weekday::Sat => 6,
        Weekday::Sun => 7,
    }
}
const APP_ID: &str = "org.gtk_rs.HelloWorld2";
const BUTTON_WIDTH: i32 = 180;
const BUTTON_HEIGHT: i32 = 120;
const PADDING: i32 = 6;
// Function to get the number of days in a month
fn get_days_in_month(year: i32, month: u32) -> u32 {
    match month {
        1 | 3 | 5 | 7 | 8 | 10 | 12 => 31,
        4 | 6 | 9 | 11 => 30,
        2 => {
            if year % 4 == 0 && (year % 100 != 0 || year % 400 == 0) {
                // Leap year
                29
            } else {
                // Non-leap year
                28
            }
        }
        _ => panic!("Invalid month"),
    }
}

fn main() -> glib::ExitCode {
    
    // Create a new application
    let app = Application::builder().application_id(APP_ID).build();

    app.connect_startup(|_| load_css());
    // Connect to "activate" signal of `app`
    app.connect_activate(build_ui);

    app.set_accels_for_action("win.close", &["<Ctrl>W"]);

    // Run the application
    app.run()
}

fn load_css() {
    // Load the CSS file and add it to the provider
    let provider = CssProvider::new();
    provider.load_from_data(include_str!("style.css"));

    // Add the provider to the default screen
    gtk::style_context_add_provider_for_display(
        &Display::default().expect("Could not connect to a display."),
        &provider,
        gtk::STYLE_PROVIDER_PRIORITY_APPLICATION,
    );
}
fn build_ui(app: &Application) {

     // Get the current local date
    let local_date = Local::now();

    // Extract the year and month
    let year = local_date.year();
    let month = local_date.month();

    // Get the number of days in the current month
    let days_in_month = get_days_in_month(year, month);

    let notebook = gtk::Notebook::new();
    let weekly = gtk::Box::builder()
        .orientation(Orientation::Vertical)
        .build();
    // Create two buttons
    let button_increase = Button::builder()
        .label("Increase")
        .margin_top(PADDING)
        .margin_bottom(PADDING)
        .margin_start(PADDING)
        .margin_end(PADDING)
        .build();
    button_increase.set_widget_name("button_increase");
    let button_decrease = Button::builder()
        .label("Decrease")
        .margin_top(PADDING)
        .margin_bottom(PADDING)
        .margin_start(PADDING)
        .margin_end(PADDING)
        .build();
    button_decrease.set_widget_name("button_decrease");
    // Reference-counted object with inner-mutability
    let number = Rc::new(Cell::new(0));

    // Connect callbacks
    // When a button is clicked, `number` and label of the other button will be changed
    button_increase.connect_clicked(clone!(@weak number, @weak button_decrease =>
        move |_| {number.set(number.get() + 1);
            button_decrease.set_label(&number.get().to_string());
    }));
    button_decrease.connect_clicked(clone!(@weak button_increase =>
        move |_| {
            number.set(number.get() - 1);
            button_increase.set_label(&number.get().to_string());
    }));
    weekly.append(&button_increase);
    weekly.append(&button_decrease);
    notebook.append_page(&weekly, Some(&gtk::Label::new(Some("Weekly"))));
    // Add buttons to `gtk_box`
    let monthly = gtk::Box::builder()
        .orientation(Orientation::Vertical)
        .css_classes(vec![String::from("hover-box")])
        .build();
    monthly.set_size_request(7*(PADDING+BUTTON_WIDTH)+PADDING, 7*(PADDING+BUTTON_HEIGHT));
    let gesture = gtk::GestureClick::new();
    gesture.connect_released(|gesture, _, _, _| {
        gesture.set_state(gtk::EventSequenceState::Claimed);
        println!("Box pressed!");
    });
    monthly.add_controller(gesture);
    // gtk_box.append(&button_increase);
    // gtk_box.append(&button_decrease);

    // Create a window
    let window = ApplicationWindow::builder()
        .application(app)
        .title("Calendar")
        // .child(&gtk_box)
        .build();
    // window.set_default_size(300, 250);
    let scrolled = gtk::ScrolledWindow::new();
    scrolled.set_policy(gtk::PolicyType::Never, gtk::PolicyType::Automatic);
    scrolled.set_size_request(7*(PADDING+BUTTON_WIDTH)+PADDING, 7*(PADDING+BUTTON_HEIGHT));
    let flowbox = gtk::FlowBox::new();
    flowbox.set_valign(gtk::Align::Start);
    flowbox.set_max_children_per_line(7);
    flowbox.set_min_children_per_line(7);
    flowbox.set_selection_mode(gtk::SelectionMode::None);
    
    let mut offset: u32 = 0;
    match NaiveDate::from_ymd_opt(year, month, 1) {
        Some(date) => {
            offset = weekday_to_u32(date.weekday());
        }
        None => {
            println!("Invalid date");
        }
    };
    for d in 1..35 {
        if d<offset {
            let button = Button::builder()
                .margin_top(PADDING)
                .margin_bottom(PADDING)
                .margin_start(PADDING)
                .margin_end(PADDING)
                .width_request(BUTTON_WIDTH)
                .height_request(BUTTON_HEIGHT)
                .build();
            button.set_visible(false);
            flowbox.append(&button);
        } else {
            let current_date = NaiveDate::from_ymd_opt(year, month, d-offset+1);
        
            // let current_date = Utc.with_ymd_and_hms(year, month, day, 0,0,0);
            let button = Button::builder()
                .margin_top(PADDING)
                .margin_bottom(PADDING)
                .margin_start(PADDING)
                .margin_end(PADDING)
                .width_request(BUTTON_WIDTH)
                .height_request(BUTTON_HEIGHT)
                .build();
            match current_date {
                Some(date) => {
                    let day_of_week = date.weekday();
                    button.set_label(&(day_of_week.to_string()));
                }
                None => {
                    println!("Invalid date");
                    break;
                }
            };
            button.connect_closure(
                "clicked",
                false,
                closure_local!(move |button: gtk::Button| {
                    // Set the label to "Hello World!" after the button has been clicked on
                    button.set_label("Hello World!");
                    
                    show_modal_window(d);
                }),
            );
            flowbox.append(&button);
        }
    }
    scrolled.set_child(Some(&flowbox));
    monthly.append(&scrolled);
    notebook.append_page(&monthly, Some(&gtk::Label::new(Some("Monthly"))));
    window.set_child(Some(&notebook));
    window.set_resizable(false);
    // Add action "close" to `window` taking no parameter
    let action_close = ActionEntry::builder("close")
    .activate(|window: &ApplicationWindow, _, _| {
        window.close();
    })
    .build();
    window.add_action_entries([action_close]);
    // window.set_default_size(500, 500);
    // Present window
    window.present();
}